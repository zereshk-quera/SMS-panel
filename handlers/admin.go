package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	database "SMS-panel/database"
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

// AdminRegisterHandler registers a new admin
// @Summary Register a new admin
// @Description Register a new admin with the provided details
// @Tags admin
// @Accept json
// @Produce json
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
	jsonFormatValidationMsg, jsonFormatErr := utils.ValidateJsonFormat(jsonBody, "firstname", "lastname", "email", "phone", "nationalid", "username", "password")
	if jsonFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: jsonFormatValidationMsg})
	}

	// check password correction
	if jsonBody["password"].(string) != (os.Getenv("ADMIN_PASSWORD")) {
		log.Println(os.Getenv("ADMIN_PASSWORD"))
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Password Isn't Correct"})
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

// This Function Used To Deactivate An Account
func DeactivateHandler(c echo.Context, db *gorm.DB) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var account models.Account

	db.First(&account, id)
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

// This Function Used To Activate An Account
func ActivateHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	var account models.Account

	db.First(&account, id)
	if account.ID == 0 {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid Account ID"})
	}

	// if account is active
	if account.IsActive {
		return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active !"})
	}

	// activate account and update database
	account.IsActive = true
	db.Save(&account)
	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active From Now"})
}

// This Function Used To Add a config
func AddConfigHandler(c echo.Context) error {
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

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
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

	return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "Configuration Added Successfuly"})
}

// This Function Used To Hide passwords with length 5 to 7 from admin in the messages
func hidePassword(message string) string {
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

// This function used to count messages for each account
func SmsReportHandler(c echo.Context) error {
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}
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

// This function used to give messages which have specific word
func SmsSearchHandler(c echo.Context) error {
	word := c.Param("word")

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	var messages []models.SMSMessage
	db.Select("id").Find(&messages)

	res := make([]string, 0)
	senders := make([]string, 0)

	// count the number of messages for each account
	for _, msg := range messages {
		var tmp models.SMSMessage
		db.First(&tmp, msg.ID)
		if strings.Contains(tmp.Message, word) {
			res = append(res, hidePassword(tmp.Message))
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

func AddBadWordHandler(c echo.Context) error {
	word := c.Param("word")

	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

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
