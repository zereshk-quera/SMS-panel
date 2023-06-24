package handlers

import (
	database "SMS-panel/database"
	"SMS-panel/models"
	"SMS-panel/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

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
	accountCreationMsg, account, accountCreationErr := utils.CreateAccount(int(user.ID), jsonBody["username"].(string), true, jsonBody["password"].(string), db)
	if accountCreationErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: accountCreationMsg})
	}

	return c.JSON(http.StatusOK, account)
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

	db, err := database.GetConnection()
	if err != nil {
		return err
	}
	findAccountMsg, account, findAccountErr := utils.Login(jsonBody["username"].(string), jsonBody["password"].(string), true, db)
	if findAccountErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: findAccountMsg})
	}

	return c.JSON(http.StatusOK, account)

}

// This Function Used To Deactivate An Account
func DeactivateHandler(c echo.Context) error {
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
	if account.ID == 1 {
		return c.JSON(http.StatusBadRequest, models.Response{ResponseCode: 400, Message: "You can't deactive super admin!"})
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

	//if account is active
	if account.IsActive {
		return c.JSON(http.StatusOK, models.Response{ResponseCode: 200, Message: "This Account is active !"})
	}

	//activate account and update database
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

	//check json format
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
		//check for digit
		if message[messageIndex] >= 48 && message[messageIndex] <= 57 {
			tmp := ""
			for messageIndex < len(message) && message[messageIndex] >= 48 && message[messageIndex] <= 57 {
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

	//count the number of messages for each account
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
	//jsonData, _ := json.Marshal(res)
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
