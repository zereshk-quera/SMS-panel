package handlers

import (
	database "SMS-panel/database"
	"SMS-panel/models"
	"SMS-panel/utils"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// This Function Used To Deactivate An Account
func DeactivateHandler(c echo.Context) error {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "id")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	idFloat, _ := jsonBody["id"].(float64)
	id := int(idFloat)
	var account models.Account

	db.First(&account, id)
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Account ID"})
	}

	//if account is deactivated before
	if !account.IsActive {
		return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account Isn't active !"})
	}

	//deactivate account and update database
	account.IsActive = false
	account.Token = ""
	db.Save(&account)
	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account Isn't active From Now"})
}

// This Function Used To Activate An Account
func ActivateHandler(c echo.Context) error {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "id")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	//id := jsonBody["id"].(int)

	idFloat, _ := jsonBody["id"].(float64)
	id := int(idFloat)
	var account models.Account

	db.First(&account, id)
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Account ID"})
	}

	//if account is active
	if account.IsActive {
		return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active !"})
	}

	//activate account and update database
	account.IsActive = true
	db.Save(&account)
	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active From Now"})

}

// This Function Used To Activate An Account
func AdminLoginHandler(c echo.Context) error {
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

	// Find account based on input username
	username := jsonBody["username"].(string)
	var account models.Account
	db, err := database.GetConnection()
	if err != nil {
		return err
	}
	db.Where("username = ?", username).First(&account)

	// Account Not Found
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Username"})
	}

	// Incorrect Password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(jsonBody["password"].(string)))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Wrond Password"})
	}

	//Account isn't an admin
	if !account.IsAdmin {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "You are not admin!"})
	}

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    account.ID,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"admin": true,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "Failed To Create Token"})
	}

	// Update Account's Token In Database
	account.Token = tokenString
	db.Save(&account)

	// Check for Cookie Existence
	hasCookie := false
	cookies := c.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "account_token" && cookie.Value == account.Token {
			hasCookie = true
			break
		}
	}

	// Create Cookie
	if !hasCookie {
		cc := &http.Cookie{
			Name:   "account_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		c.SetCookie(cc)
		c.SetCookie(&http.Cookie{Name: "account_token", MaxAge: -1})

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

// This Function Used To Add a config
func AddConfigHandler(c echo.Context) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "name", "value")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	var conf models.Configuration
	conf.Name = jsonBody["name"].(string)
	conf.Value = jsonBody["value"].(float64)

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	createdConf := db.Create(&conf)
	if createdConf.Error != nil {
		return c.JSON(http.StatusInternalServerError, conf)
	}

	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "Configuration Added Successfuly"})
}

// This Function Used To Registeraion of a New Admin
func AdminRegisterHandler(c echo.Context) error {
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

	//check password correction
	if jsonBody["password"].(string) != (os.Getenv("ADMIN_PASSWORD")) {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Password Isn't Correct"})
	}

	// Create User Object
	var user models.User
	user.FirstName = jsonBody["firstname"].(string)
	user.LastName = jsonBody["lastname"].(string)
	user.Email = jsonBody["email"].(string)
	user.Phone = jsonBody["phone"].(string)
	user.NationalID = jsonBody["nationalid"].(string)

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

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
	db.Where("phone = ?", user.Phone).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Phone Number has already been registered"})
	}

	// Is Input Email Address Unique or Not
	db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Email Address has already been registered"})
	}

	// Is Input National ID Unique or Not
	db.Where("national_id = ?", user.NationalID).First(&existingUser)
	if existingUser.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt National ID has already been registered"})
	}

	// Is Input Username Unique or Not
	var existingAccount models.Account
	db.Where("username = ?", jsonBody["username"].(string)).First(&existingAccount)
	if existingAccount.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Inupt Username has already been registered"})
	}

	// Insert User Object Into Database
	createdUser := db.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "User Cration Failed"})
	}

	// Instantiating Account Object
	var account models.Account
	account.UserID = user.ID
	account.Username = jsonBody["username"].(string)
	account.Budget = 0
	account.IsAdmin = true

	hash, err := bcrypt.GenerateFromPassword([]byte(jsonBody["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 422, Message: "Failed to Hashing Password"})
	}

	account.Password = string(hash)

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    account.ID,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"admin": true,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "Failed To Create Token"})
	}
	account.Token = tokenString

	// Insert Account Object Into Database
	createdAccount := db.Create(&account)
	if createdAccount.Error != nil {
		return c.JSON(http.StatusInternalServerError, account)
	}

	cc := &http.Cookie{
		Name:   "account_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	c.SetCookie(cc)
	c.SetCookie(&http.Cookie{Name: "account_token", MaxAge: -1})

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

// This Function Used To Hide passwords with length 5 to 7 from admin in the messages
func hidePassword(message string) string {
	messageIndex := 0
	res := ""
	for messageIndex < len(message) {
		//check for digit
		if message[messageIndex] >= 48 && message[messageIndex] <= 57 {
			tmp := ""
			for message[messageIndex] >= 48 && message[messageIndex] <= 57 {
				tmp += string(message[messageIndex])
				messageIndex++
			}
			//check length of number
			if len(tmp) >= 5 && len(tmp) <= 7 {
				for i := 0; i < len(tmp); i++ {
					res += "_"
				}
			} else {
				res += tmp
			}
		} else {
			res += string(message[messageIndex])
			messageIndex++
		}
	}
	return res
}

// This function used to count messages for each account
func SmsReportHandler(c echo.Context) error {

	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}
	//get all accounts
	var accounts []models.Account

	db.Select("id").Find(&accounts)

	//fill map with account ids as keys
	var accountIDs = make(map[string]int, 0)
	for _, account := range accounts {
		accountIDs["Account "+strconv.FormatUint(uint64(account.ID), 10)] = 0
	}

	//get all sms messages
	var messages []models.SMSMessage
	db.Select("id").Find(&messages)

	//count the number of messages for each account
	for _, msg := range messages {
		var tmp models.SMSMessage
		db.First(&tmp, msg.ID)
		accountIDs["Account "+strconv.FormatUint(uint64(tmp.AccountID), 10)]++
	}
	return c.JSON(http.StatusOK, accountIDs)
}

// This function used to give messages which have specific word
func SmsSearchHandler(c echo.Context) error {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "word")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	word := jsonBody["word"].(string)

	var messages []models.SMSMessage
	db.Select("id").Find(&messages)

	res := make([]string, 0)

	//count the number of messages for each account
	for _, msg := range messages {
		var tmp models.SMSMessage
		db.First(&tmp, msg.ID)
		if strings.Contains(tmp.Message, word) {
			res = append(res, hidePassword(tmp.Message))
		}
	}
	ret := make(map[string]string)
	i := 1
	for _, v := range res {
		t := "Message "
		t += strconv.Itoa(i)
		ret[t] = v
		i++
	}
	//jsonData, _ := json.Marshal(res)
	return c.JSON(http.StatusOK, ret)
}

func AddBadWordHandler(c echo.Context) error {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	//check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "word")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	word := jsonBody["word"].(string)
	reg := utils.GenerateRegex(word)

	var bw models.Bad_Word
	bw.Word = word
	bw.Regex = reg
	createdBW := db.Create(&bw)
	if createdBW.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "Bad Word Cration Failed"})
	}
	return c.JSON(http.StatusOK, bw)
}
