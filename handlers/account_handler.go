package handlers

import (
	database "SMS-panel/database"
	"SMS-panel/models"
	"SMS-panel/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var SECRET = "s89ut8cn4u3bghyn75gy38ghm9g3mgc85g9m" ///should be in env file !!!!!!!!!!

func RegisterHandler(c echo.Context) error {
	//Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//Create User Object
	var user models.User
	user.FirstName = jsonBody["firstname"].(string)
	user.LastName = jsonBody["lastname"].(string)
	user.Email = jsonBody["email"].(string)
	user.Phone = jsonBody["phone"].(string)
	user.NationalID = jsonBody["nationalid"].(string)

	//Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	//Check FirstName Validation
	if len(strings.TrimSpace(user.FirstName)) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "First Name can't be empty"})
	}

	//Check LastName Validation
	if len(strings.TrimSpace(user.LastName)) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Last Name can't be empty"})
	}

	//Check Phone Number Validation
	if !utils.ValidatePhone(user.Phone) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Phone Number"})
	}

	//Check Email Validation
	if !utils.ValidateEmail(user.Email) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Email Address"})
	}

	//Check NationalID Validation
	if !utils.ValidateNationalID(user.NationalID) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid National ID"})
	}

	//Is Input Phone Number Unique or Not
	var existingUser models.User
	db.Where("phone = ?", user.Phone).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Phone Number has already been registered"})
	}

	//Is Input Email Address Unique or Not
	db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Email Address has already been registered"})
	}

	//Is Input Email Address Unique or Not
	db.Where("nationalid = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt National ID has already been registered"})
	}

	//Is Input Username Unique or Not
	var existingAccount models.Account
	db.Where("username = ?", jsonBody["username"].(string)).First(&existingAccount)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Username has already been registered"})
	}

	//Insert User Object Into Database
	createdUser := db.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "User Cration Failed"})
	}
	userID := user.ID

	//instantiating Account Object
	var account models.Account
	account.UserID = userID
	account.Username = jsonBody["username"].(string)
	account.Budget = 0

	hash, err := bcrypt.GenerateFromPassword([]byte(jsonBody["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 422, Message: "Failed to Hashing Password"})
	}

	account.Password = string(hash)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "Failed To Create Token"})
	}
	account.Password = tokenString

	//Insert Account Object Into Database
	createdAccount := db.Create(&account)
	if createdAccount.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "Account Cration Failed"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = account.Token
	cookie.Expires = time.Now().Add(time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, account)
}

func LoginHandler(c echo.Context) error {
	//Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	username := jsonBody["username"].(string)
	var account models.Account
	db, err := database.GetConnection()
	db.Where("username = ?", username).First(&account)

	//Account Not Found
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Username"})
	}

	//Incorrect Password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(jsonBody["password"].(string)))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Wrond Password"})
	}

	//Generate Token
	var user models.User
	db.Where("id = ?", account.UserID).First(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "Failed To Create Token"})
	}
	account.Token = tokenString

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = account.Token
	cookie.Expires = time.Now().Add(time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, account)
}

func BudgetAmountHandler(c echo.Context) error {
	account := c.Get("account")
	budget := account.(models.Account).Budget
	res := struct {
		amount int
	}{
		amount: int(budget),
	}
	return c.JSON(http.StatusOK, res)
}
