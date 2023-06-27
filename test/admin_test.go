package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
