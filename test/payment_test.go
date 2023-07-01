package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRequestHandler(t *testing.T) {
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
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)

	t.Run("InvalidJSON", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts/payment/request", strings.NewReader(`{ "invalid": "json" }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.PaymentRequestHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Invalid JSON", response.Message)
	})

	t.Run("MissingFee", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts/payment/request", strings.NewReader(`{}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("account", account)

		err := handlers.PaymentRequestHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "Input Json doesn't include fee", response.Message)
	})
	t.Run("SuccessfulRequest", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts/payment/request", strings.NewReader(`{ "fee": 100000 }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("account", account)

		err := handlers.PaymentRequestHandler(c, db)
		assert.NoError(t, err)
		log.Println(rec.Body.String())

		assert.Equal(t, http.StatusOK, rec.Code)

		var response handlers.RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response.PaymentUrl)
	})

	t.Run("FeeUnderTheMinimum", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts/payment/request", strings.NewReader(`{ "fee": 100 }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("account", account)

		err := handlers.PaymentRequestHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 400, int(response.ResponseCode))
		assert.Equal(t, "fee must not be under 1000", response.Message)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		err := db.Delete(&user).Error
		assert.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts/payment/request", strings.NewReader(`{ "fee": 10000 }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("account", account)

		err = handlers.PaymentRequestHandler(c, db)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 404, int(response.ResponseCode))
		assert.Equal(t, "User Not Founded", response.Message)
	})
}
