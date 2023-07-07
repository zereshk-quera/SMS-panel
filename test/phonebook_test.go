package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreatePhoneBook(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	defer utils.CloseTestDatabase(db)
	assert.NoError(t, err)

	handler := handlers.NewPhonebookHandler(db)
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

	e := echo.New()

	t.Run("CreatePhoneBookSuccess", func(t *testing.T) {
		requestBody := handlers.PhoneBookRequest{
			AccountID: account.ID,
			Name:      "Test Phone Book",
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBook(ctx)
		log.Println(rec.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response handlers.PhoneBookResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Assert the response
		assert.Equal(t, requestBody.AccountID, response.AccountID)
		assert.Equal(t, requestBody.Name, response.Name)
	})
	t.Run("CreatePhoneBookNameMissing", func(t *testing.T) {
		requestBody := handlers.PhoneBookRequest{
			AccountID: account.ID,
			Name:      "",
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBook(ctx)
		log.Println(rec.Body.String())

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"error":"Name is required"}`, strings.TrimSpace(rec.Body.String()))
	})
}

func TestGetAllPhoneBooks(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	defer utils.CloseTestDatabase(db)
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
		Username: "testuser",
		Budget:   10,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)
	handler := handlers.NewPhonebookHandler(db)

	phoneBooks := []models.PhoneBook{
		{AccountID: account.ID, Name: "Phone Book 1"},
		{AccountID: account.ID, Name: "Phone Book 2"},
	}
	err = db.Create(&phoneBooks).Error
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/account/"+fmt.Sprint(account.ID)+"/phone-books/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("accountID")
	ctx.SetParamValues(fmt.Sprint(account.ID))
	ctx.Set("account", account)

	err = handler.GetAllPhoneBooks(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []handlers.PhoneBookResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, uint(1), response[0].ID)
	assert.Equal(t, uint(account.ID), response[0].AccountID)
	assert.Equal(t, "Phone Book 1", response[0].Name)
	assert.Equal(t, uint(2), response[1].ID)
	assert.Equal(t, uint(account.ID), response[1].AccountID)
	assert.Equal(t, "Phone Book 2", response[1].Name)
}

func TestReadPhoneBook(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	defer utils.CloseTestDatabase(db)
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
		Username: "testuser",
		Budget:   10,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)

	handler := handlers.NewPhonebookHandler(db)
	accountID := account.ID

	phoneBook := models.PhoneBook{
		AccountID: account.ID,
		Name:      "Phone Book 1",
	}
	err = db.Create(&phoneBook).Error
	assert.NoError(t, err)
	phoneBookID := phoneBook.ID

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/account/"+fmt.Sprint(accountID)+"/phone-books/"+fmt.Sprint(phoneBookID), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("accountID", "phoneBookID")
	ctx.SetParamValues(fmt.Sprint(accountID), fmt.Sprint(phoneBookID))
	ctx.Set("account", account)

	err = handler.ReadPhoneBook(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response handlers.PhoneBookResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(phoneBookID), response.ID)
	assert.Equal(t, uint(accountID), response.AccountID)
	assert.Equal(t, "Phone Book 1", response.Name)
}

func TestUpdatePhoneBook(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	defer utils.CloseTestDatabase(db)
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
		Username: "testuser",
		Budget:   10,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	handler := handlers.NewPhonebookHandler(db)
	accountID := account.ID

	phoneBook := models.PhoneBook{
		AccountID: accountID,
		Name:      "Phone Book 1",
	}
	err = db.Create(&phoneBook).Error
	assert.NoError(t, err)
	phoneBookID := phoneBook.ID

	e := echo.New()

	updateData := map[string]interface{}{
		"name": "Updated Phone Book",
	}
	updateJSON, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPut, "/account/"+fmt.Sprint(accountID)+"/phone-books/"+fmt.Sprint(phoneBookID), bytes.NewReader(updateJSON))
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("accountID", "phoneBookID")
	ctx.SetParamValues(fmt.Sprint(accountID), fmt.Sprint(phoneBookID))
	ctx.Set("account", account)

	err = handler.UpdatePhoneBook(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response handlers.PhoneBookResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(phoneBookID), response.ID)
	assert.Equal(t, uint(accountID), response.AccountID)
	assert.Equal(t, "Updated Phone Book", response.Name)

	var updatedPhoneBook models.PhoneBook
	err = db.Where("id = ?", phoneBook.ID).First(&updatedPhoneBook).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Phone Book", updatedPhoneBook.Name)
}

func TestDeletePhoneBook(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	defer utils.CloseTestDatabase(db)
	assert.NoError(t, err)

	handler := handlers.NewPhonebookHandler(db)

	account := models.Account{
		UserID:   1,
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
		Name:      "Test Phone Book",
	}
	err = db.Create(&phoneBook).Error
	assert.NoError(t, err)

	e := echo.New()

	t.Run("DeletePhoneBookSuccess", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/account/"+strconv.Itoa(int(account.ID))+"/phone-books/"+strconv.Itoa(int(phoneBook.ID)), nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("accountID", "phoneBookID")
		ctx.SetParamValues(strconv.Itoa(int(account.ID)), strconv.Itoa(int(phoneBook.ID)))
		ctx.Set("account", account)

		err := handler.DeletePhoneBook(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"Phone book deleted\"\n", rec.Body.String())

		var deletedPhoneBook models.PhoneBook
		result := db.First(&deletedPhoneBook, phoneBook.ID)
		assert.Error(t, result.Error)
		assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
	})

	t.Run("DeletePhoneBookNotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/account/"+strconv.Itoa(int(account.ID))+"/phone-books/999", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("accountID", "phoneBookID")
		ctx.SetParamValues(strconv.Itoa(int(account.ID)), "999")
		ctx.Set("account", account)

		err := handler.DeletePhoneBook(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Phone book not found\"", strings.TrimSpace(rec.Body.String()))
	})
}
