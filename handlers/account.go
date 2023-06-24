package handlers

import (
	"encoding/json"
	"net/http"

	database "SMS-panel/database"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
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
func RegisterHandler(c echo.Context) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "firstname", "lastname", "email", "phone", "nationalid", "username", "password")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	//check user validation
	userFormatValidationMsg, user, userFormatErr := utils.ValidateUser(jsonBody)
	if userFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: userFormatValidationMsg})
	}

	//check unique
	userUniqueMsg, userUniqueErr := utils.CheckUnique(user, jsonBody["username"].(string), db)
	if userUniqueErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: userUniqueMsg})
	}

	// Insert User Object Into Database
	createdUser := db.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "User Cration Failed"})
	}

	//create account
	accountCreationMsg, account, accountCreationErr := utils.CreateAccount(int(user.ID), jsonBody["username"].(string), false, jsonBody["password"].(string), db)
	if accountCreationErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: accountCreationMsg})
	}

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
func LoginHandler(c echo.Context) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "username", "password")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	//get database connection
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	//find account based on username and check password correction
	findAccountMsg, account, findAccountErr := utils.Login(jsonBody["username"].(string), jsonBody["password"].(string), false, db)
	if findAccountErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: findAccountMsg})
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
func BudgetAmountHandler(c echo.Context) error {
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
