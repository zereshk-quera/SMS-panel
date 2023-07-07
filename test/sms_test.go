package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSendSingleSMSHandler(t *testing.T) {
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
		Budget:   10,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)
	phoneBook := models.PhoneBook{
		AccountID: account.ID,
		Name:      "Test",
	}
	err = db.Create(&phoneBook).Error
	assert.NoError(t, err)
	phonebooknumber := models.PhoneBookNumber{
		PhoneBookID: phoneBook.ID,
		Username:    "test",
		Name:        "test",
		Phone:       "09376304339",
	}
	err = db.Create(&phonebooknumber).Error
	assert.NoError(t, err)
	config := models.Configuration{
		Name:  "single sms",
		Value: 100,
	}
	err = db.Create(&config).Error
	assert.NoError(t, err)

	t.Run("InvalidRequestPayload", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "invalid": "json" `)) // Modify the JSON to make it invalid
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "Invalid request payload", response.Message)
	})

	t.Run("InvalidPhoneNumber", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "phone_number": "123" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "Invalid phone number", response.Message)
	})

	t.Run("SenderNumberNotFound", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "SenderNumber": "123456", "PhoneNumber": "1234567890" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)

		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "Sender number not found!", response.Message)
	})

	t.Run("InsufficientBudget", func(t *testing.T) {
		e := echo.New()
		sender_number := models.SenderNumber{
			Number:    "123456789",
			IsDefault: true,
		}
		err = db.Create(&sender_number).Error
		assert.NoError(t, err)
		user_number := models.UserNumbers{
			UserID:      user.ID,
			NumberID:    sender_number.ID,
			EndDate:     time.Now().AddDate(1, 0, 0),
			StartDate:   time.Now(),
			IsAvailable: true,
			Number:      sender_number,
			User:        user,
		}
		err = db.Create(&user_number).Error
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "senderNumbers": "123456789", "phone_number": "09376304339","message":"hello","username":"test" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)
		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusForbidden, rec.Code)

		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusForbidden, response.Code)
		assert.Equal(t, "Insufficient budget", response.Message)
	})

	t.Run("UsernameNotFound", func(t *testing.T) {
		newBudget := 200
		err = db.Model(&account).Update("budget", newBudget).Error
		assert.NoError(t, err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "senderNumbers": "123456789", "phone_number": "09376304339","message":"hello","username":"username" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "Username not found", response.Message)
	})

	t.Run("PhoneNumberNotFound", func(t *testing.T) {
		newBudget := 200
		err = db.Model(&account).Update("budget", newBudget).Error
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "senderNumbers": "123456789", "phone_number": "09376304331","message":"hello" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "Phone number not found", response.Message)
	})

	t.Run("SuccessfulSMS", func(t *testing.T) {
		newBudget := 200
		err = db.Model(&account).Update("budget", newBudget).Error
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/single-sms", strings.NewReader(`{ "senderNumbers": "123456789", "phone_number": "09376304339","message":"hello","username":"test" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.SendSingleSMSHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response handlers.SendSMSResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "SMS sent successfully", response.Message)
	})
}

func TestPeriodicSendSMSHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
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
		Username: "test",
		Budget:   10000,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)
	phoneBook := models.PhoneBook{
		Name:      "Test Phone Book",
		AccountID: account.ID,
	}
	err = db.Create(&phoneBook).Error
	assert.NoError(t, err)

	phoneBookNumber1 := models.PhoneBookNumber{
		Phone:       "09376304339",
		PhoneBookID: phoneBook.ID,
		PhoneBook:   phoneBook,
		Username:    "test",
		Name:        "test",
	}
	err = db.Create(&phoneBookNumber1).Error
	assert.NoError(t, err)
	phoneBookNumber2 := models.PhoneBookNumber{
		Phone:       "0987654321",
		PhoneBookID: phoneBook.ID,
		PhoneBook:   phoneBook,
		Username:    "test2",
		Name:        "test2",
	}
	err = db.Create(&phoneBookNumber2).Error
	assert.NoError(t, err)
	sender_number := models.SenderNumber{
		Number:    "123456789",
		IsDefault: true,
	}
	err = db.Create(&sender_number).Error
	assert.NoError(t, err)
	user_number := models.UserNumbers{
		UserID:      user.ID,
		NumberID:    sender_number.ID,
		EndDate:     time.Now().AddDate(1, 0, 0),
		StartDate:   time.Now(),
		IsAvailable: true,
		Number:      sender_number,
		User:        user,
	}
	err = db.Create(&user_number).Error
	assert.NoError(t, err)
	config := models.Configuration{
		Name:  "single sms",
		Value: 100,
	}
	err = db.Create(&config).Error
	assert.NoError(t, err)
	config2 := models.Configuration{
		Name:  "group sms",
		Value: 100,
	}
	err = db.Create(&config2).Error
	assert.NoError(t, err)

	defer utils.CloseTestDatabase(db)

	t.Run("ValidRequest", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", strings.NewReader(`{
			"username": "test",
			"senderNumbers": "123456789",
			"phone": "09376304339",
			"message": "Hello",
			"schedule": "09:00",
			"interval": "daily"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "SMS scheduled successfully", rec.Body.String())
	})

	t.Run("InvalidRequestPayload", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", strings.NewReader("invalid request payload"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request payload", rec.Body.String())
	})

	t.Run("InvalidScheduleTimeFormat", func(t *testing.T) {
		e := echo.New()
		reqBody := handlers.SendSMSRequestPeriodic{
			Username:     "test",
			SenderNumber: "123456789",
			Phone:        "",
			Message:      "Hello",
			Schedule:     "invalid",
			Interval:     "daily",
			PhoneBookID:  "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid schedule time format", rec.Body.String())
	})

	t.Run("RecipientNotProvided", func(t *testing.T) {
		e := echo.New()
		reqBody := handlers.SendSMSRequestPeriodic{
			Username:     "",
			SenderNumber: "123456789",
			Phone:        "",
			Message:      "Hello",
			Schedule:     "09:00",
			Interval:     "daily",
			PhoneBookID:  "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Recipient not provided", rec.Body.String())
	})

	t.Run("SenderNumberNotFound", func(t *testing.T) {
		e := echo.New()
		reqBody := handlers.SendSMSRequestPeriodic{
			Username:     "test",
			SenderNumber: "unknown",
			Phone:        "",
			Message:      "Hello",
			Schedule:     "09:00",
			Interval:     "daily",
			PhoneBookID:  "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response handlers.ErrorResponseSingle
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "Sender number not found!", response.Message)
	})

	t.Run("InsufficientBudget", func(t *testing.T) {
		newBudget := 0
		err = db.Model(&account).Update("budget", newBudget).Error
		assert.NoError(t, err)
		e := echo.New()
		reqBody := handlers.SendSMSRequestPeriodic{
			Username:     "test",
			SenderNumber: "123456789",
			Phone:        "09376304339",
			Message:      "Hello",
			Schedule:     "09:00",
			Interval:     "daily",
			PhoneBookID:  "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Insufficient budget", rec.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		newBudget := 1000
		err = db.Model(&account).Update("budget", newBudget).Error
		e := echo.New()
		reqBody := handlers.SendSMSRequestPeriodic{
			Username:     "",
			SenderNumber: "123456789",
			Phone:        "09376304339",
			Message:      "Hello",
			Schedule:     "09:00",
			Interval:     "daily",
			PhoneBookID:  "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/sms/periodic-sms", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("account", account)

		err := handlers.PeriodicSendSMSHandler(c, db)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "SMS scheduled successfully", rec.Body.String())
	})
}
