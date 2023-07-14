package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminRegistrationRequest struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	NationalID string `json:"nationalid"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
type ConfigurationRequest struct {
	Name  string  `json:"name" example:"config_name"`
	Value float64 `json:"value" example:"42.0"`
}

// AdminRegisterHandler registers a new admin
// @Summary Register a new admin
// @Description Register a new admin with the provided details
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "User Token"
// @Param adminRegistrationRequest body AdminRegistrationRequest true "Admin registration details"
// @Success 200 {object} AccountResponse "Admin account created successfully"
// @Failure 422 {object} ErrorResponse "Invalid JSON"
// @Failure 422 {object} ErrorResponse "JSON format validation failed"
// @Failure 422 {object} ErrorResponse "Input password isn't correct"
// @Failure 422 {object} ErrorResponse "User validation failed"
// @Failure 422 {object} ErrorResponse "Username already exists"
// @Failure 500 {object} ErrorResponse "User creation failed"
// @Failure 502 {object} ErrorResponse "Can't connect to the database"
// @Router /admin/register [post]
func AdminRegisterHandler(c echo.Context, db *gorm.DB) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "firstname", "lastname", "email", "phone", "nationalid", "username", "password", "code")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// check code correction
	adminCode, err := utils.GetAdminCode()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "Failed to get admin code"})
	}

	if jsonBody["code"].(string) != adminCode {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Code Isn't Correct"})
	}

	// check user validation
	userFormatValidationMsg, user, userFormatErr := utils.ValidateUser(jsonBody)
	if userFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: userFormatValidationMsg})
	}

	// check unique
	userUniqueMsg, userUniqueErr := utils.CheckUnique(user, jsonBody["username"].(string), db)
	if userUniqueErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: userUniqueMsg})
	}

	// Insert User Object Into Database
	createdUser := db.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "User Creation Failed"})
	}

	// create account
	accountCreationMsg, account, accountCreationErr := utils.CreateAccount(int(user.ID), jsonBody["username"].(string), true, jsonBody["password"].(string), db)
	if accountCreationErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: accountCreationMsg})
	}

	return c.JSON(http.StatusOK, account)
}

// AdminLoginHandler logs in an admin user.
// @Summary Admin login
// @Description Logs in an admin user.
// @Tags admin
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login Request Body"
// @Success 200 {object} AccountResponse
// @Failure 422 {object} ErrorResponseRegisterLogin
// @Router /admin/login [post]
func AdminLoginHandler(c echo.Context, db *gorm.DB) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "username", "password")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	findAccountMsg, account, findAccountErr := utils.Login(jsonBody["username"].(string), jsonBody["password"].(string), true, db)
	if findAccountErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: findAccountMsg})
	}

	return c.JSON(http.StatusOK, account)
}

// DeactivateHandler deactivates an account.
// @Summary Deactivate Account
// @Description Deactivates the specified account.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 422 {object} models.Response
// @Router /admin/deactivate/{id} [patch]
// @Router /admin/deactivate/{id} [patch]
func DeactivateHandler(c echo.Context, db *gorm.DB) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var account models.Account

	db.Where("id = ?", id).First(&account)
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Account ID"})
	}
	if account.ID == 1 {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "You can't deactive super admin!"})
	}

	// if account is deactivated before
	if !account.IsActive {
		return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account Isn't active !"})
	}

	// deactivate account and update database
	account.IsActive = false
	account.Token = ""
	db.Save(&account)
	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account Isn't active From Now"})
}

// ActivateHandler activates an account.
// @Summary Activate Account
// @Description Activates the specified account.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 422 {object} models.Response
// @Router /admin/activate/{id} [patch]
func ActivateHandler(c echo.Context, db *gorm.DB) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var account models.Account

	db.Where("id = ?", id).First(&account)
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Account ID"})
	}

	// if account is active
	if account.IsActive {
		return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active!"})
	}

	// activate account and update database
	account.IsActive = true
	db.Save(&account)
	log.Println("Account Activated:", account)

	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active From Now"})
}

// AddConfigHandler creates a new configuration entry.
// @Summary Create Configuration
// @Description Create a new configuration entry
// @Tags admin
// @Accept json
// @Produce json
// @Param config body ConfigurationRequest true "Configuration object to be added"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 422 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/add-config [post]
func AddConfigHandler(c echo.Context, db *gorm.DB) error {
	// Read Request Body
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	// check json format
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "name", "value")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	var existingConfig models.Configuration
	db.Where("name = ?", jsonBody["name"].(string)).First(&existingConfig)
	if existingConfig.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "There is a config with input name in database"})
	}

	var conf models.Configuration
	conf.Name = jsonBody["name"].(string)
	conf.Value = jsonBody["value"].(float64)

	createdConf := db.Create(&conf)
	if createdConf.Error != nil {
		return c.JSON(http.StatusInternalServerError, conf)
	}

	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "Configuration Added Successfully"})
}

// This Function Used To Hide passwords with length 5 to 7 from admin in the messages
func HidePassword(message string) string {
	messageIndex := 0
	res := ""
	for messageIndex < len(message) {
		// check for digit
		if message[messageIndex] >= 48 && message[messageIndex] <= 57 {
			tmp := ""
			for messageIndex < len(message) && message[messageIndex] >= 48 && message[messageIndex] <= 57 {
				tmp += string(message[messageIndex])
				messageIndex++
			}
			// check length of number
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

// SmsReportHandler retrieves the SMS report.
// @Summary Get SMS Report
// @Description Retrieve the SMS report with the number of messages per account
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} map[string]int
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/sms-report [get]
func SmsReportHandler(c echo.Context, db *gorm.DB) error {
	// get all accounts
	var accounts []models.Account

	db.Select("id").Find(&accounts)

	// fill map with account ids as keys
	accountIDs := make(map[string]int, 0)
	for _, account := range accounts {
		accountIDs["Account "+strconv.FormatUint(uint64(account.ID), 10)] = 0
	}

	// get all sms messages
	var messages []models.SMSMessage
	db.Select("id").Find(&messages)

	// count the number of messages for each account
	for _, msg := range messages {
		var tmp models.SMSMessage
		db.First(&tmp, msg.ID)
		accountIDs["Account "+strconv.FormatUint(uint64(tmp.AccountID), 10)]++
	}
	return c.JSON(http.StatusOK, accountIDs)
}

// SmsSearchHandler searches for SMS messages containing a specific word.
// @Summary Search SMS Messages
// @Description Search for SMS messages containing a specific word
// @Tags admin
// @Param word path string true "Word to search for in SMS messages"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} map[string]string
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/search/{word} [get]
func SmsSearchHandler(c echo.Context, db *gorm.DB) error {
	word := c.Param("word")

	var messages []models.SMSMessage
	db.Select("id").Find(&messages)

	res := make([]string, 0)
	senders := make([]string, 0)

	// count the number of messages for each account
	for _, msg := range messages {
		var tmp models.SMSMessage
		db.First(&tmp, msg.ID)
		if strings.Contains(tmp.Message, word) {
			res = append(res, HidePassword(tmp.Message))
			senders = append(senders, tmp.Sender)
		}
	}
	ret := make(map[string]string)
	counter := 1
	for i, v := range res {
		t := fmt.Sprintf("%d. %s ", counter, senders[i])
		ret[t] = v
		counter++
	}
	// jsonData, _ := json.Marshal(res)
	return c.JSON(http.StatusOK, ret)
}

// AddBadWordHandler
// @Summary Add a bad word
// @Description Add a new bad word to the database
// @Tags admin
// @Accept json
// @Produce json
// @Param word path string true "Word to add as a bad word"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} models.Bad_Word
// @Failure 422 {object} models.Response
// @Failure 500 {object} models.Response
// @Failure 502 {object} models.Response
// @Router /admin/add-bad-word/{word} [post]
func AddBadWordHandler(c echo.Context, db *gorm.DB) error {
	word := c.Param("word")

	var existingWord models.Bad_Word
	db.Where("word = ?", word).First(&existingWord)
	if existingWord.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "This word added before as a bad word in database"})
	}

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
