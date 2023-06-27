package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
