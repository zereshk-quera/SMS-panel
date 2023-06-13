package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	database "SMS-panel/database"
	"SMS-panel/handlers"
	"SMS-panel/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	account     models.Account
	phoneBookID uint
)

func TestMain(m *testing.M) {
	err := database.Connect()
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}
	/*
		err = cleanupTestData()
		if err != nil {
			panic("failed to cleanup test data: " + err.Error())
		}
	*/

	err = createTestData()
	if err != nil {
		panic("failed to create test data: " + err.Error())
	}

	code := m.Run()

	err = cleanupTestData()
	if err != nil {
		panic("failed to cleanup test data: " + err.Error())
	}

	os.Exit(code)
}

func createTestData() error {
	user := models.User{
		FirstName:  "John",
		LastName:   "Doe",
		Phone:      "123456789",
		Email:      "john.doe@example.com",
		NationalID: "1234567890",
	}
	db, err := database.GetConnection()
	if err != nil {
		return err
	}
	err = db.Create(&user).Error
	if err != nil {
		return err
	}

	account = models.Account{
		UserID:   user.ID,
		Username: "johndoe",
		Budget:   1000,
		Password: "password",
	}
	err = db.Create(&account).Error
	if err != nil {
		return err
	}

	phoneBook := models.PhoneBook{
		AccountID: account.ID,
		Name:      "John Doe",
	}
	err = db.Create(&phoneBook).Error
	if err != nil {
		return err
	}

	phoneBookID = phoneBook.ID

	return nil
}

func cleanupTestData() error {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	err = db.Exec("DELETE FROM phone_books WHERE account_id = ?", account.ID).Error
	if err != nil {
		return err
	}

	err = db.Exec("DELETE FROM accounts WHERE id = ?", account.ID).Error
	if err != nil {
		return err
	}
	err = db.Exec("DELETE FROM users WHERE id = ?", account.UserID).Error
	if err != nil {
		return err
	}

	return nil
}

func TestCreatePhoneBook(t *testing.T) {
	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		phoneBookReq := handlers.PhoneBookRequest{
			AccountID: account.ID,
			Name:      "John Doe",
		}

		reqBody, err := json.Marshal(phoneBookReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/phonebook", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handlers.CreatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var phoneBookRes handlers.PhoneBookResponse
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBookRes)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookReq.AccountID, phoneBookRes.AccountID)
		assert.Equal(t, phoneBookReq.Name, phoneBookRes.Name)
	})

	t.Run("MissingName", func(t *testing.T) {
		phoneBookReq := handlers.PhoneBookRequest{
			AccountID: account.ID,
		}

		reqBody, err := json.Marshal(phoneBookReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/phonebook", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handlers.CreatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorRes map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &errorRes)
		assert.NoError(t, err)

		assert.Equal(t, "Name is required", errorRes["error"])
	})
}

func TestGetAllPhoneBooks(t *testing.T) {
	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/phonebooks/"+fmt.Sprint(account.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID")
		c.SetParamValues(fmt.Sprint(account.ID))

		err := handlers.GetAllPhoneBooks(c)
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
	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, phoneBookID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(phoneBookID))

		err := handlers.ReadPhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var phoneBook models.PhoneBook
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBook)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookID, phoneBook.ID)
		assert.Equal(t, account.ID, phoneBook.AccountID)
		assert.Equal(t, "John Doe", phoneBook.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, 99999), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), "99999")

		err := handlers.ReadPhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorMessage string
		err = json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		assert.NoError(t, err)

		assert.Equal(t, "Phonebook not found", errorMessage)
	})
}

func TestUpdatePhoneBook(t *testing.T) {
	e := echo.New()

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

		err = handlers.UpdatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var phoneBookRes handlers.PhoneBookResponse
		err = json.Unmarshal(rec.Body.Bytes(), &phoneBookRes)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookReq.AccountID, phoneBookRes.AccountID)
		assert.Equal(t, phoneBookID, phoneBookRes.ID)
		assert.Equal(t, phoneBookReq.Name, phoneBookRes.Name)
	})

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

		err = handlers.UpdatePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorMessage string
		err = json.Unmarshal(rec.Body.Bytes(), &errorMessage)
		assert.NoError(t, err)

		assert.Equal(t, "Phonebook not found", errorMessage)
	})
}

func TestDeletePhoneBook(t *testing.T) {
	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, phoneBookID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(phoneBookID))

		err := handlers.DeletePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponseBody := "\"Phone book deleted\""
		actualResponseBody := strings.TrimSpace(rec.Body.String())
		assert.Equal(t, expectedResponseBody, actualResponseBody)
	})

	t.Run("NotFound", func(t *testing.T) {
		nonExistentPhoneBookID := phoneBookID + 100
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/account/%d/phone-books/%d", account.ID, nonExistentPhoneBookID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("accountID", "phoneBookID")
		c.SetParamValues(fmt.Sprint(account.ID), fmt.Sprint(nonExistentPhoneBookID))

		err := handlers.DeletePhoneBook(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		expectedResponseBody := "\"Phone book not found\""
		actualResponseBody := strings.TrimSpace(rec.Body.String())
		assert.Equal(t, expectedResponseBody, actualResponseBody)
	})
}
