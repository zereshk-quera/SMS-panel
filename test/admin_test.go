package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAdminRegisterHandler(t *testing.T) {
	e := echo.New()

	t.Run("ValidRequest", func(t *testing.T) {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		requestBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "0817762590",
			"username":   "johndoe",
			"password":   adminPassword,
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/admin/register", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		err = handlers.AdminRegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Account
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "johndoe", response.Username)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/register", bytes.NewReader([]byte("invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		err = handlers.AdminRegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(rec.Code))
		assert.Equal(t, "Invalid JSON", response.Message)
	})

	t.Run("InvalidUserFormat", func(t *testing.T) {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		requestBody := map[string]interface{}{
			"firstname":  "",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "1234567890",
			"nationalid": "123456789",
			"username":   "johndoe",
			"password":   adminPassword,
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/admin/register", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		err = handlers.AdminRegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "First Name can't be empty", response.Message)
	})

	t.Run("IncorrectPassword", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "0817762590",
			"username":   "johndoe",
			"password":   "incorrectpassword",
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/admin/register", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		err = handlers.AdminRegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Input Password Isn't Correct", response.Message)
	})
	t.Run("NonUniquePhoneNumber", func(t *testing.T) {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		requestBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "1234567890", // Non-unique phone number
			"nationalid": "0817762590",
			"username":   "johndoe",
			"password":   adminPassword,
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/admin/register", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)
		user := models.User{
			FirstName:  "john",
			LastName:   "doe",
			Phone:      "1234567890",
			Email:      "existing@example.com",
			NationalID: "987654321",
		}
		err = db.Create(&user).Error
		assert.NoError(t, err)

		err = handlers.AdminRegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Invalid Phone Number", response.Message)
	})
}

func TestAdminLoginHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	e := echo.New()
	user := models.User{FirstName: "testuser", LastName: "testuser", Phone: "09376304339", Email: "amir@gmail.com", NationalID: "0265670578"}
	db.Create(&user)
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	account := models.Account{UserID: user.ID, Username: "admin", Password: string(hash), Token: "testtoken", IsActive: true, IsAdmin: true}
	db.Create(&account)

	t.Run("ValidRequest", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"username": "admin",
			"password": "admin123",
		}
		assert.NoError(t, err)
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err = handlers.AdminLoginHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Account
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewReader([]byte("invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err = handlers.AdminLoginHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Invalid JSON", response.Message)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"username": "admin",
			"password": "wrongpassword",
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = handlers.AdminLoginHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Wrong Password", response.Message)
	})
}

func TestDeactivateHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)
	user := models.User{
		FirstName:  "john",
		LastName:   "doe",
		Phone:      "09376304339",
		Email:      "test@gmail.com",
		NationalID: "123456789",
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	account := models.Account{
		UserID:   user.ID,
		Username: "testuser",
		Budget:   0,
		Password: "password",
		IsActive: true,
		IsAdmin:  true,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)
	testAccount := models.Account{
		UserID:   user.ID,
		Username: "testuser2",
		Budget:   0,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	err = db.Create(&testAccount).Error
	assert.NoError(t, err)

	t.Run("ValidAccountID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/admin/deactivate/%d", testAccount.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", testAccount.ID))

		err := handlers.DeactivateHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(response.ResponseCode))
		assert.Equal(t, "This Account Isn't active From Now", response.Message)

		updatedAccount := models.Account{}
		db.First(&updatedAccount, testAccount.ID)
		assert.False(t, updatedAccount.IsActive)
	})

	t.Run("InvalidAccountID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/admin/deactivate/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handlers.DeactivateHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Invalid Account ID", response.Message)
	})

	t.Run("SuperAdminAccount", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/admin/deactivate/%d", account.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", account.ID))

		err := handlers.DeactivateHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 400, int(response.ResponseCode))
		assert.Equal(t, "You can't deactive super admin!", response.Message)

		updatedAccount := models.Account{}
		db.First(&updatedAccount, account.ID)
		assert.True(t, updatedAccount.IsActive)
	})
}

func TestActivateHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)

	user := models.User{
		FirstName:  "john",
		LastName:   "doe",
		Phone:      "09376304339",
		Email:      "test@gmail.com",
		NationalID: "123456789",
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)

	account := models.Account{
		UserID:   user.ID,
		Username: "testuser",
		Budget:   0,
		Password: "password",
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)

	t.Run("ValidAccountID", func(t *testing.T) {
		err = db.Model(&account).Update("IsActive", false).Error
		assert.NoError(t, err)

		log.Println("account activate status", account.IsActive)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/admin/activate/%d", account.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", account.ID))

		err := handlers.ActivateHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(response.ResponseCode))
		assert.Equal(t, "This Account is active From Now", response.Message)

		updatedAccount := models.Account{}
		err = db.First(&updatedAccount, account.ID).Error
		assert.NoError(t, err)

		assert.True(t, updatedAccount.IsActive)
	})

	t.Run("InvalidAccountID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/admin/activate/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handlers.ActivateHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Invalid Account ID", response.Message)
	})

	t.Run("AlreadyActiveAccount", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/admin/activate/%d", account.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", account.ID))

		account.IsActive = true
		db.Save(&account)

		err := handlers.ActivateHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 200, int(response.ResponseCode))
		assert.Equal(t, "This Account is active!", response.Message)
	})
}

func TestAddConfigHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)

	e := echo.New()

	t.Run("ValidRequest", func(t *testing.T) {
		reqBody := `{"name": "config1", "value": 10}`
		req := httptest.NewRequest(http.MethodPost, "/admin/add-config", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err = handlers.AddConfigHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		expectedResponse := models.Response{
			ResponseCode: 200,
			Message:      "Configuration Added Successfully",
		}
		assert.Equal(t, expectedResponse, response)

		var conf models.Configuration
		db.First(&conf, "name = ?", "config1")
		assert.Equal(t, "config1", conf.Name)
		assert.Equal(t, 10.0, conf.Value)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		reqBody := `{"name": "config1"}`
		req := httptest.NewRequest(http.MethodPost, "/admin/add-config", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		// Call the handler function
		err = handlers.AddConfigHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		expectedResponse := models.Response{
			ResponseCode: 422,
			Message:      "Input Json doesn't include value",
		}
		assert.Equal(t, expectedResponse, response)

		var count int64
		db.Model(&models.Configuration{}).Count(&count)
		assert.Equal(t, int64(1), count) // set this one beacuse one will created in first test
	})

	t.Run("DuplicateName", func(t *testing.T) {
		conf := models.Configuration{
			Name:  "config1",
			Value: 10.0,
		}
		db.Create(&conf)

		reqBody := `{"name": "config1", "value": 20}`
		req := httptest.NewRequest(http.MethodPost, "/admin/add-config", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.Set("db", db)

		err = handlers.AddConfigHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		expectedResponse := models.Response{
			ResponseCode: 422,
			Message:      "There is a config with input name in database",
		}
		assert.Equal(t, expectedResponse, response)

		var count int64
		db.Model(&models.Configuration{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestHidePassword(t *testing.T) {
	t.Run("NoDigits", func(t *testing.T) {
		message := "No digits in this message"
		expectedResult := "No digits in this message"
		result := handlers.HidePassword(message)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("HideShortNumber", func(t *testing.T) {
		message := "12345 is a short number"
		expectedResult := "_____ is a short number"
		result := handlers.HidePassword(message)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("HideMediumNumber", func(t *testing.T) {
		message := "1234567 is a medium number"
		expectedResult := "_______ is a medium number"
		result := handlers.HidePassword(message)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("HideLongNumber", func(t *testing.T) {
		message := "12345678 is a long number"
		expectedResult := "12345678 is a long number"
		result := handlers.HidePassword(message)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("MultipleNumbers", func(t *testing.T) {
		message := "12345 is a short number, and 987654321 is a long number"
		expectedResult := "_____ is a short number, and 987654321 is a long number"
		result := handlers.HidePassword(message)
		assert.Equal(t, expectedResult, result)
	})
}

func TestSmsReportHandler(t *testing.T) {
	t.Run("NoAccounts", func(t *testing.T) {
		// Create a test database and defer its closure
		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		// Create an Echo instance and set up the request
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/sms-report", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the handler
		err = handlers.SmsReportHandler(c, db)

		// Assert the handler's response
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var accountIDs map[string]int
		err = json.Unmarshal(rec.Body.Bytes(), &accountIDs)
		assert.NoError(t, err)

		// Assert the result
		assert.Empty(t, accountIDs)
	})

	t.Run("WithAccountsAndMessages", func(t *testing.T) {
		// Create a test database and defer its closure
		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)
		user := models.User{
			FirstName:  "testuser",
			LastName:   "testuser",
			Phone:      "09376304339",
			Email:      "amir@gmail.com",
			NationalID: "123456789",
		}
		err = db.Create(&user).Error
		assert.NoError(t, err)
		accounts := []models.Account{
			{ID: 1, UserID: 1, Username: "testuser1", Budget: 0, Password: "test", Token: "test", IsActive: true, IsAdmin: false},
			{ID: 2, UserID: 1, Username: "testuser2", Budget: 0, Password: "test", Token: "test", IsActive: true, IsAdmin: false},
			{ID: 3, UserID: 1, Username: "testuser3", Budget: 0, Password: "test", Token: "test", IsActive: true, IsAdmin: false},
		}
		err = db.Create(&accounts).Error
		assert.NoError(t, err)

		messages := []models.SMSMessage{
			{ID: 1, AccountID: 1, Sender: "123456789", Recipient: "12345678", Message: "message", DeliveryReport: "done"},
			{ID: 2, AccountID: 1, Sender: "123456789", Recipient: "12345678", Message: "message", DeliveryReport: "done"},
			{ID: 3, AccountID: 2, Sender: "123456789", Recipient: "12345678", Message: "message", DeliveryReport: "done"},
			{ID: 4, AccountID: 2, Sender: "123456789", Recipient: "12345678", Message: "message", DeliveryReport: "done"},
			{ID: 5, AccountID: 3, Sender: "123456789", Recipient: "12345678", Message: "message", DeliveryReport: "done"},
			{ID: 6, AccountID: 3, Sender: "123456789", Recipient: "12345678", Message: "message", DeliveryReport: "done"},
		}
		err = db.Create(&messages).Error
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/sms-report", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = handlers.SmsReportHandler(c, db)

		// Assert the handler's response
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var accountIDs map[string]int
		err = json.Unmarshal(rec.Body.Bytes(), &accountIDs)
		assert.NoError(t, err)

		// Assert the result
		assert.Equal(t, 3, len(accountIDs))
		assert.Equal(t, 2, accountIDs["Account 1"])
		assert.Equal(t, 2, accountIDs["Account 2"])
		assert.Equal(t, 2, accountIDs["Account 3"])
	})
}

func TestSmsSearchHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)
	user := models.User{
		FirstName:  "testuser",
		LastName:   "testuser",
		Phone:      "09376304339",
		Email:      "amir@gmail.com",
		NationalID: "123456789",
	}
	err = db.Create(&user).Error
	assert.NoError(t, err)
	accounts := []models.Account{
		{ID: 1, UserID: 1, Username: "testuser1", Budget: 0, Password: "test", Token: "test", IsActive: true, IsAdmin: false},
	}
	err = db.Create(&accounts).Error
	assert.NoError(t, err)

	t.Run("WordFound", func(t *testing.T) {
		messages := []models.SMSMessage{
			{ID: 1, AccountID: 1, Sender: "123456789", Recipient: "12345678", Message: "message1", DeliveryReport: "done"},
			{ID: 2, AccountID: 1, Sender: "123456789", Recipient: "12345678", Message: "message2", DeliveryReport: "done"},
		}
		err = db.Create(&messages).Error
		assert.NoError(t, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/search/{word}", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("word")
		c.SetParamValues("message")

		err = handlers.SmsSearchHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		log.Println(rec.Body.String())
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, 2)
		assert.Equal(t, "message1", response["1. 123456789 "])
		assert.Equal(t, "message2", response["2. 123456789 "])
	})

	t.Run("WordNotFound", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/search/{word}", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("word")
		c.SetParamValues("test")

		err = handlers.SmsSearchHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response)
	})

	t.Run("MessagesFoundWithPassword", func(t *testing.T) {
		messages := []models.SMSMessage{
			{ID: 3, AccountID: 1, Sender: "123456789", Recipient: "12345678", Message: "test1234560", DeliveryReport: "done"},
			{ID: 4, AccountID: 1, Sender: "123456789", Recipient: "12345678", Message: "test2345678", DeliveryReport: "done"},
		}
		err = db.Create(&messages).Error
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/search/{word}", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("word")
		c.SetParamValues("test")

		err = handlers.SmsSearchHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(response))
		assert.Equal(t, "test_______", response["1. 123456789 "])
		assert.Equal(t, "test_______", response["2. 123456789 "])
	})
}

func TestAddBadWordHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)

	t.Run("SuccessfulCreation", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/add-bad-word/:word", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("word")
		c.SetParamValues("example")

		err := handlers.AddBadWordHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var bw models.Bad_Word
		err = json.Unmarshal(rec.Body.Bytes(), &bw)
		assert.NoError(t, err)

		assert.Equal(t, "example", bw.Word)
		assert.NotEmpty(t, bw.Regex)
	})

	t.Run("FailureExistingWord", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/add-bad-word/:word", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetParamNames("word")
		c.SetParamValues("example")

		err = handlers.AddBadWordHandler(c, db)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "This word added before as a bad word in database", response.Message)
	})
}
