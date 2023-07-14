package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
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

type RentNumberRequest struct {
	SenderNumber              string `json:"senderNumber"`
	SubscriptionNumberPackage string `json:"SubscriptionNumberPackage"`
}

type BuyNumberRequest struct {
	SenderNumber string `json:"senderNumber"`
}

type BudgetAmountResponse struct {
	Amount int `json:"amount"`
}

type SenderNumbersResponse struct {
	Numbers []string `json:"numbers"`
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
func RegisterHandler(c echo.Context, dbConn *gorm.DB) error {
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

	// check user validation
	userFormatValidationMsg, user, userFormatErr := utils.ValidateUser(jsonBody)
	if userFormatErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: userFormatValidationMsg})
	}

	// check unique
	userUniqueMsg, userUniqueErr := utils.CheckUnique(user, jsonBody["username"].(string), dbConn)
	if userUniqueErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: userUniqueMsg})
	}

	// Insert User Object Into Database
	createdUser := dbConn.Create(&user)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "User Creation Failed"})
	}

	// create account
	accountCreationMsg, account, accountCreationErr := utils.CreateAccount(int(user.ID), jsonBody["username"].(string), false, jsonBody["password"].(string), dbConn)
	if accountCreationErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: accountCreationMsg})
	}

	return c.JSON(http.StatusCreated, account)
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
func LoginHandler(c echo.Context, db *gorm.DB) error {
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

	// find account based on username and check password correction
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
// @Param Authorization header string true "User Token"
// @Success 200 {object} BudgetAmountResponse
// @Failure 401 {string} string
// @Router /accounts/budget [get]
func BudgetAmountHandler(c echo.Context) error {
	// Recieve Account Object
	account := c.Get("account")

	account = account.(models.Account)
	budget := int(account.(models.Account).Budget)

	// Create Result Object
	res := BudgetAmountResponse{
		Amount: budget,
	}
	return c.JSON(http.StatusOK, res)
}

// GetAllSenderNumbersHandler retrieves All sender numbers available for the account
// @Summary Get All sender numbers
// @Description retrieves All sender numbers available for the account
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} SenderNumbersResponse
// @Failure 401 {string} string
// @Router /accounts/sender-numbers	 [get]
func GetAllSenderNumbersHandler(c echo.Context, db *gorm.DB) error {
	account := c.Get("account").(models.Account)

	var senderNumbersObjects []models.SenderNumber

	err := db.Model(&models.SenderNumber{}).
		Select("sender_numbers.number").
		Joins("LEFT JOIN user_numbers ON sender_numbers.id = user_numbers.number_id").
		Where("sender_numbers.is_default=true or (user_numbers.user_id = ? and user_numbers.is_available=true)",
			account.UserID).
		Scan(&senderNumbersObjects).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}
	var senderNumbers []string
	for _, n := range senderNumbersObjects {
		senderNumbers = append(senderNumbers, n.Number)
	}

	return c.JSON(http.StatusOK, SenderNumbersResponse{Numbers: senderNumbers})
}

// GetAllSenderNumbersForSaleHandler retrieves All sender numbers available for sale
// @Summary Get All sender numbers for sale
// @Description retrieves All sender numbers available for sale
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} SenderNumbersResponse
// @Failure 401 {string} string
// @Router /accounts/sender-numbers/sale	 [get]
func GetAllSenderNumbersForSaleHandler(c echo.Context, db *gorm.DB) error {
	var senderNumbersObjects []models.SenderNumber

	err := db.Model(&models.SenderNumber{}).
		Select("number").
		Where("is_default = false AND is_exclusive=false").
		Scan(&senderNumbersObjects).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}
	var senderNumbers []string
	for _, n := range senderNumbersObjects {
		senderNumbers = append(senderNumbers, n.Number)
	}

	return c.JSON(http.StatusOK, SenderNumbersResponse{Numbers: senderNumbers})
}

// @Summary Rent number
// @Description Rent available number for this account
// @Tags users
// @Accept json
// @Produce json
// @Param body body RentNumberRequest true "Get sender number and subscription package."
// @Success 200 {object} models.Response
// @Failure 204 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/rent-number [post]
func RentNumberHandler(c echo.Context, db *gorm.DB) error {
	account := c.Get("account").(models.Account)
	body := RentNumberRequest{}
	ctx := c.Request().Context()

	if err := c.Bind(&body); err != nil {
		errResponse := ErrorResponse{
			Message: "Invalid request payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	// Check if sender number is available for this user
	var senderNumbersObject models.SenderNumber

	err := db.WithContext(ctx).WithContext(ctx).Model(&models.SenderNumber{}).
		Select("sender_numbers.id", "sender_numbers.number").
		Joins("LEFT JOIN user_numbers ON sender_numbers.id = user_numbers.number_id").
		Where(
			"sender_numbers.is_default=false and sender_numbers.is_exclusive=false and sender_numbers.number = ?",
			body.SenderNumber).
		First(&senderNumbersObject).Error
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Sender number not found!"})
	}

	// get the subscription number package
	var subPackage models.SubscriptionNumberPackage
	err = db.WithContext(ctx).
		Where("title = ?", body.SubscriptionNumberPackage).
		First(&subPackage).Error
	if err != nil {
		errorResponse := ErrorResponse{Message: "Subscription package does not exist."}
		return c.JSON(http.StatusNotFound, errorResponse)
	}
	haveAccountBudget := utils.DoesAcountHaveBudget(
		account.Budget, subPackage.Price,
	)
	if !haveAccountBudget {
		errorResponse := ErrorResponse{Message: "You don't have enough budget!"}
		return c.JSON(http.StatusNotFound, errorResponse)
	}

	var subscriptionNumberPackage SubscriptionNumberPackageInterface
	startDate := time.Now()
	if subPackage.Title == "1 Month" {
		subscriptionNumberPackage = &OneMonthSubscriptionNumberPackage{StartDate: startDate}
	} else if subPackage.Title == "2 Month" {
		subscriptionNumberPackage = &TwoMonthSubscriptionNumberPackage{StartDate: startDate}
	}

	// Save to database
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	startDate, endDate := subscriptionNumberPackage.GetTimePeriod()
	userNumberObject := models.UserNumbers{
		UserID:                account.UserID,
		NumberID:              senderNumbersObject.ID,
		StartDate:             startDate,
		EndDate:               endDate,
		IsAvailable:           true,
		SubscriptionPackageID: subPackage.ID,
	}
	if err = tx.Create(&userNumberObject).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	// Update senderNumber
	err = tx.Model(&models.SenderNumber{}).Where("id = ?", senderNumbersObject.ID).
		Update("is_exclusive", true).
		Error
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	// Update account budget
	account.Budget -= subPackage.Price
	if err = tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, models.Response{
		ResponseCode: http.StatusOK, Message: "success",
	})
}

// @Summary Buy number
// @Description Buy available number for this account
// @Tags users
// @Accept json
// @Produce json
// @Param body body BuyNumberRequest true "Get sender number and subscription package."
// @Success 200 {object} models.Response
// @Failure 204 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/buy-number [post]
func BuyNumberHandler(c echo.Context, db *gorm.DB) error {
	account := c.Get("account").(models.Account)
	body := BuyNumberRequest{}
	ctx := c.Request().Context()

	if err := c.Bind(&body); err != nil {
		errResponse := ErrorResponse{
			Message: "Invalid request payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	// Check if sender number is available for this user
	var senderNumbersObject models.SenderNumber

	err := db.WithContext(ctx).WithContext(ctx).Model(&models.SenderNumber{}).
		Select("sender_numbers.id", "sender_numbers.number").
		Joins("LEFT JOIN user_numbers ON sender_numbers.id = user_numbers.number_id").
		Where(
			"sender_numbers.is_default=false and sender_numbers.is_exclusive=false and sender_numbers.number = ?",
			body.SenderNumber).
		First(&senderNumbersObject).Error
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Sender number not found!"})
	}

	// get the subscription number package
	var subPackage models.SubscriptionNumberPackage
	err = db.WithContext(ctx).First(
		&subPackage, utils.SUBSCRIOPTION_PACKAGE_BUY_ID,
	).Error
	if err != nil {
		errorResponse := ErrorResponse{Message: "Subscription package does not exist."}
		return c.JSON(http.StatusNotFound, errorResponse)
	}
	haveAccountBudget := utils.DoesAcountHaveBudget(
		account.Budget, subPackage.Price,
	)
	if !haveAccountBudget {
		errorResponse := ErrorResponse{Message: "You don't have enough budget!"}
		return c.JSON(http.StatusNotFound, errorResponse)
	}

	// Save to database
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	userNumberObject := models.UserNumbers{
		UserID:                account.UserID,
		NumberID:              senderNumbersObject.ID,
		StartDate:             time.Now(),
		IsAvailable:           true,
		SubscriptionPackageID: subPackage.ID,
	}
	if err = tx.Create(&userNumberObject).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	// Update senderNumber
	err = tx.Model(&models.SenderNumber{}).Where("id = ?", senderNumbersObject.ID).
		Update("is_exclusive", true).
		Error
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	// Update account budget
	account.Budget -= subPackage.Price
	if err = tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error"})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, models.Response{
		ResponseCode: http.StatusOK, Message: "success",
	})
}
