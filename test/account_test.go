package test

import (
	"SMS-panel/handlers"
	"SMS-panel/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// @Router /accounts/register
func TestRegisterHandler(t *testing.T) {

	t.Run("Success Register", func(t *testing.T) {

		registerUserCreate := handlers.UserCreateRequest{
			FirstName:  "Rick",
			LastName:   "Sanchez",
			Email:      "RickSanchez@morty.com",
			Phone:      "09123456789",
			NationalID: "0369734971",
			Username:   "ricksanchez",
			Password:   "123Rick123",
		}

		reqBody, err := json.Marshal(registerUserCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/accounts/register", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handlers.RegisterHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// create hash for password
		hash, err := bcrypt.GenerateFromPassword([]byte(registerUserCreate.Password), bcrypt.DefaultCost)
		assert.NoError(t, err)
		registerAccount := handlers.AccountResponse{
			Username: registerUserCreate.Username,
			Password: string(hash),
			Budget:   0,
		}
		var accountResponse handlers.AccountResponse
		err = json.Unmarshal(rec.Body.Bytes(), &accountResponse)
		assert.NoError(t, err)

		assert.Equal(t, registerAccount.Username, accountResponse.Username)
		assert.Equal(t, registerAccount.Password, accountResponse.Password)
		assert.Equal(t, registerAccount.Budget, accountResponse.Budget)

	})

	t.Run("Stupid User without FirstName input", func(t *testing.T) {

		withoutFirstNameUserCreate := handlers.UserCreateRequest{
			//FirstName: "Rick",
			LastName:   "Sanchez",
			Email:      "RickSanchez@morty.com",
			Phone:      "09123456789",
			NationalID: "0369734971",
			Username:   "ricksanchez",
			Password:   "123Rick123",
		}

		reqBody, err := json.Marshal(withoutFirstNameUserCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/accounts/register", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handlers.RegisterHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var responseMessage models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &responseMessage)
		assert.NoError(t, err)

		assert.Equal(t, 422, responseMessage.ResponseCode)
		assert.Equal(t, "Input Json doesn't include firstname", responseMessage.Message)

	})
	t.Run("Stupid User without LastName input", func(t *testing.T) {

		withoutLastNameUserCreate := handlers.UserCreateRequest{
			FirstName: "ًRick",
			//LastName: "Sanchez",
			Email:      "RickSanchez@morty.com",
			Phone:      "09123456789",
			NationalID: "0369734971",
			Username:   "ricksanchez",
			Password:   "123Rick123",
		}

		reqBody, err := json.Marshal(withoutLastNameUserCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/accounts/register", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handlers.RegisterHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var responseMessage models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &responseMessage)
		assert.NoError(t, err)

		assert.Equal(t, 422, responseMessage.ResponseCode)
		assert.Equal(t, "Input Json doesn't include lastname", responseMessage.Message)

	})

	t.Run("Stupid User without Email input", func(t *testing.T) {

		withoutEmailUserCreate := handlers.UserCreateRequest{
			FirstName: "ًRick",
			LastName:  "Sanchez",
			//Email:      "RickSanchez@morty.com",
			Phone:      "09123456789",
			NationalID: "0369734971",
			Username:   "ricksanchez",
			Password:   "123Rick123",
		}

		reqBody, err := json.Marshal(withoutEmailUserCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/accounts/register", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handlers.RegisterHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var responseMessage models.Response
		err = json.Unmarshal(rec.Body.Bytes(), &responseMessage)
		assert.NoError(t, err)

		assert.Equal(t, 422, responseMessage.ResponseCode)
		assert.Equal(t, "Input Json doesn't include email", responseMessage.Message)

	})
}
