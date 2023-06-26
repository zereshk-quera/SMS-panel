package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	database "SMS-panel/database"
	"SMS-panel/handlers"
	"SMS-panel/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	e = echo.New()
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}

	phonebookHandler = handlers.NewPhonebookHandler(db)
}

func TestCreatePhoneBook(t *testing.T) {
	t.Skip("skipping for now")
	t.Run("Success", func(t *testing.T) {
		phoneBookReq := handlers.PhoneBookRequest{
			AccountID: account.ID,
			Name:      "John Doe",
		}

		reqBody, err := json.Marshal(phoneBookReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = phonebookHandler.CreatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var phoneBookRes handlers.PhoneBookResponse
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBookRes)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookReq.AccountID, phoneBookRes.AccountID)
		assert.Equal(t, phoneBookReq.Name, phoneBookRes.Name)
	})
	t.Skip("skipping for now")
	t.Run("MissingName", func(t *testing.T) {
		phoneBookReq := handlers.PhoneBookRequest{
			AccountID: account.ID,
		}

		reqBody, err := json.Marshal(phoneBookReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = phonebookHandler.CreatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorRes map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errorRes)
		assert.NoError(t, err)

		assert.Equal(t, "Name is required", errorRes["error"])
	})
}

func TestGetAllPhoneBooks(t *testing.T) {
	t.Skip("skipping for now")
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/account/phone-books/"+fmt.Sprint(account.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID")
		c.SetParamValues(fmt.Sprint(account.ID))

		err := phonebookHandler.GetAllPhoneBooks(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var phoneBooks []models.PhoneBook
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBooks)
		assert.NoError(t, err)
		assert.Equal(t, phoneBookID, phoneBooks[0].ID)
		assert.Equal(t, account.ID, phoneBooks[0].AccountID)
		assert.Equal(t, "John Doe", phoneBooks[0].Name)
	})
}

func TestReadPhoneBook(t *testing.T) {
	t.Skip("skipping for now")
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, phoneBookID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(phoneBookID))

		err := phonebookHandler.ReadPhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var phoneBook models.PhoneBook
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBook)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookID, phoneBook.ID)
		assert.Equal(t, account.ID, phoneBook.AccountID)
		assert.Equal(t, "John Doe", phoneBook.Name)
	})
	t.Skip("skipping for now")
	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, 99999), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), "99999")

		err := phonebookHandler.ReadPhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorMessage string
		err = json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		assert.NoError(t, err)

		assert.Equal(t, "Phonebook not found", errorMessage)
	})
}

func TestUpdatePhoneBook(t *testing.T) {
	t.Skip("skipping for now")
	t.Run("Success", func(t *testing.T) {
		phoneBookReq := handlers.PhoneBookRequest{
			AccountID: account.ID,
			Name:      "Updated Name",
		}

		reqBody, err := json.Marshal(phoneBookReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, phoneBookID), bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(phoneBookID))

		err = phonebookHandler.UpdatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var phoneBookRes handlers.PhoneBookResponse
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBookRes)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookReq.AccountID, phoneBookRes.AccountID)
		assert.Equal(t, phoneBookID, phoneBookRes.ID)
		assert.Equal(t, phoneBookReq.Name, phoneBookRes.Name)
	})
	t.Skip("skipping for now")
	t.Run("NotFound", func(t *testing.T) {
		phoneBookReq := handlers.PhoneBookRequest{
			AccountID: account.ID,
			Name:      "Updated Name",
		}

		reqBody, err := json.Marshal(phoneBookReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, 99999), bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), "99999")

		err = phonebookHandler.UpdatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorMessage string
		err = json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		assert.NoError(t, err)

		assert.Equal(t, "Phonebook not found", errorMessage)
	})
}

func TestDeletePhoneBook(t *testing.T) {
	t.Skip("skipping for now")
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, phoneBookID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(phoneBookID))

		err := phonebookHandler.DeletePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponseBody := "\"Phone book deleted\""
		actualResponseBody := strings.TrimSpace(rec.Body.String())
		assert.Equal(t, expectedResponseBody, actualResponseBody)
	})
	t.Skip("skipping for now")

	t.Run("NotFound", func(t *testing.T) {
		nonExistentPhoneBookID := phoneBookID + 100
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, nonExistentPhoneBookID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(nonExistentPhoneBookID))

		err := phonebookHandler.DeletePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		expectedResponseBody := "\"Phone book not found\""
		actualResponseBody := strings.TrimSpace(rec.Body.String())
		assert.Equal(t, expectedResponseBody, actualResponseBody)
	})
}
