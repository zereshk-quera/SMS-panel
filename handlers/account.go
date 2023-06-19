package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// define this structs for swagger docs
type AccountResponse struct {
	ID       uint   `json:"ID"`
	UserID   uint   `json:"UserID"`
	Username string `json:"Username"`
	Budget   int    `json:"Budget"`
	Password string `json:"Password"`
	Token    string `json:"Token"`
	IsActive bool   `json:"IsActive"`
}
type ErrorResponseRegisterLogin struct {
	ResponseCode int    `json:"responsecode"`
	Message      string `json:"message"`
}
type UserCreateRequest struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	NationalID string `json:"nationalid"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type BudgetAmountResponse struct {
	Amount int `json:"amount"`
}

type AccountHandler struct {
	db *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param body body UserCreateRequest true "User registration details"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} ErrorResponseRegisterLogin
// @Failure 422 {object} ErrorResponseRegisterLogin
// @Failure 500 {object} ErrorResponseRegisterLogin
// @Router /accounts/register [post]
func (a AccountHandler) RegisterHandler(c echo.Context) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// check json format
	if _, ok := jsonBody["firstname"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include firstname"})
	}
	if _, ok := jsonBody["lastname"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include lastname"})
	}
	if _, ok := jsonBody["email"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include email"})
	}
	if _, ok := jsonBody["phone"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include phone"})
	}
	if _, ok := jsonBody["nationalid"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include nationalid"})
	}
	if _, ok := jsonBody["username"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include username"})
	}
	if _, ok := jsonBody["password"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include password"})
	}

	// Create User Object
	var user models.User
	user.FirstName = jsonBody["firstname"].(string)
	user.LastName = jsonBody["lastname"].(string)
	user.Email = jsonBody["email"].(string)
	user.Phone = jsonBody["phone"].(string)
	user.NationalID = jsonBody["nationalid"].(string)

	// Check FirstName Validation
	if len(strings.TrimSpace(user.FirstName)) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "First Name can't be empty"})
	}

	// Check LastName Validation
	if len(strings.TrimSpace(user.LastName)) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Last Name can't be empty"})
	}

	// Check Phone Number Validation
	if !utils.ValidatePhone(user.Phone) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Phone Number"})
	}

	// Check Email Validation
	if !utils.ValidateEmail(user.Email) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Email Address"})
	}

	// Check NationalID Validation
	if !utils.ValidateNationalID(user.NationalID) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid National ID"})
	}

	// Is Input Phone Number Unique or Not
	var existingUser models.User
	a.db.Where("phone = ?", user.Phone).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Phone Number has already been registered"})
	}

	// Is Input Email Address Unique or Not
	a.db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Email Address has already been registered"})
	}

	// Is Input National ID Unique or Not
	a.db.Where("national_id = ?", user.NationalID).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt National ID has already been registered"})
	}

	// Is Input Username Unique or Not
	var existingAccount models.Account
	a.db.Where("username = ?", jsonBody["username"].(string)).First(&existingAccount)
	if existingAccount.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Username has already been registered"})
	}

	// Insert User Object Into Database
	createdUser := a.db.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "User Cration Failed"})
	}

	// Instantiating Account Object
	var account models.Account
	account.UserID = user.ID
	account.Username = jsonBody["username"].(string)
	account.Budget = 0

	hash, err := bcrypt.GenerateFromPassword([]byte(jsonBody["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 422, Message: "Failed to Hashing Password"})
	}

	account.Password = string(hash)

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "Failed To Create Token"})
	}
	account.Token = tokenString

	// Insert Account Object Into Database
	createdAccount := a.db.Create(&account)
	if createdAccount.Error != nil {
		return c.JSON(http.StatusInternalServerError, account)
	}

	// Create Cookie
	cookie := &http.Cookie{
		Name:     "account_token",
		Value:    account.Token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, account)
}

// LoginHandler handles user login
// @Summary User login
// @Description Login with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login request body"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} ErrorResponseRegisterLogin
// @Failure 422 {object} ErrorResponseRegisterLogin
// @Router  /accounts/login [post]
func (a AccountHandler) LoginHandler(c echo.Context) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// Check json format
	if _, ok := jsonBody["username"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include username"})
	}
	if _, ok := jsonBody["password"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include password"})
	}

	// Find account based on input username
	username := jsonBody["username"].(string)
	var account models.Account

	a.db.Where("username = ?", username).First(&account)

	// Account Not Found
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Username"})
	}

	// Incorrect Password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(jsonBody["password"].(string)))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Wrond Password"})
	}

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "Failed To Create Token"})
	}

	// Update Account's Token In Database
	account.Token = tokenString
	a.db.Save(&account)

	// Check for Cookie Existence
	hasCookie := false
	cookies := c.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "account_token" {
			hasCookie = true
			break
		}
	}

	// Create Cookie
	if !hasCookie {
		cookie := &http.Cookie{
			Name:     "account_token",
			Value:    account.Token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
		}
		c.SetCookie(cookie)
	}

	return c.JSON(http.StatusOK, account)
}

// BudgetAmountHandler retrieves the budget amount for the logged-in user
// @Summary Get budget amount
// @Description Get the budget amount for the logged-in user
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} BudgetAmountResponse
// @Failure 401 {string} string
// @Router /accounts/budget	 [get]
func (a AccountHandler) BudgetAmountHandler(c echo.Context) error {
	// Recieve Account Object
	account := c.Get("account")

	account = account.(models.Account)
	budget := int(account.(models.Account).Budget)

	// Create Result Object
	res := struct {
		Amount int `json:"amount"`
	}{
		Amount: budget,
	}
	return c.JSON(http.StatusOK, res)
}
