package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler(t *testing.T) {
	e := echo.New()

	t.Run("ValidRequest", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "0817762590",
			"username":   "johndoe",
			"password":   "password",
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		err = handlers.RegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response models.Account
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "johndoe", response.Username)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := handlers.RegisterHandler(c, nil)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(rec.Code))
		assert.Equal(t, "Invalid JSON", response.Message)
	})

	t.Run("InvalidUserFormat", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"firstname":  "",
			"lastname":   "doe",
			"email":      "johndoe@example.com",
			"phone":      "1234567890",
			"nationalid": "123456789",
			"username":   "johndoe",
			"password":   "password",
		}
		jsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		db, err := utils.CreateTestDatabase()
		assert.NoError(t, err)
		defer utils.CloseTestDatabase(db)

		err = handlers.RegisterHandler(c, db)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, 422, int(response.ResponseCode))
		assert.Equal(t, "First Name can't be empty", response.Message)
	})
}

func TestLoginHandler(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)

	user := models.User{FirstName: "testuser", LastName: "testuser", Phone: "09376304339", Email: "amir@gmail.com", NationalID: "0265670578"}
	db.Create(&user)
	hash, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	account := models.Account{UserID: user.ID, Username: "testuser", Password: string(hash), Token: "testtoken"}
	db.Create(&account)

	testCases := map[string]struct {
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		"ValidLogin": {
			requestBody:    `{"username": "testuser", "password": "test123"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		"InvalidJSON": {
			requestBody:    `{"username": "testuser"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"responsecode": 422, "message": "Input Json doesn't include password"}`,
		},
		"InvalidCredentials": {
			requestBody:    `{"username": "testuser", "password": "wrongpassword"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"responsecode": 422, "message": "Wrong Password"}`,
		},
	}

	// Run the test cases
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/accounts/login", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := handlers.LoginHandler(ctx, db)
			assert.NoError(t, err, "Expected no error in LoginHandler")

			assert.Equal(t, testCase.expectedStatus, rec.Code, "Expected response status code to match")

			if testCase.expectedBody != "" {
				assert.JSONEq(t, testCase.expectedBody, rec.Body.String(), "Expected response body to match")
			}
		})
	}
}
