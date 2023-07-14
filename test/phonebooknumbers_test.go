package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"SMS-panel/handlers"
	"SMS-panel/models"
	"SMS-panel/utils"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePhoneBookNumber(t *testing.T) {
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
		Budget:   0,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}
	err = db.Create(&account).Error
	assert.NoError(t, err)
	phoneBook := models.PhoneBook{
		AccountID: account.ID,
		Name:      "Phone Book 1",
	}
	err = db.Create(&phoneBook).Error
	assert.NoError(t, err)

	e := echo.New()

	t.Run("CreatePhoneBookNumberSuccess", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "John Doe",
			Phone:       "09376304339",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: phoneBook.ID,
			PhoneBook:   phoneBook,
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response models.PhoneBookNumber
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		log.Println(response)
		assert.Equal(t, requestBody.Name, response.Name)
		assert.Equal(t, requestBody.Phone, response.Phone)
		assert.Equal(t, requestBody.Prefix, response.Prefix)
		assert.Equal(t, requestBody.Username, response.Username)
	})

	t.Run("CreatePhoneBookNumberMissingName", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "", // Empty Name
			Phone:       "09376304339",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: phoneBook.ID,
			PhoneBook:   phoneBook,
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Name is required", response["error"])
	})
	t.Run("CreatePhoneBookNumberMissingPhoneBook", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "john doe",
			Phone:       "09376304339",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: 0,
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Phone book is required", response["error"])
	})

	t.Run("CreatePhoneBookNumberMissingPhone", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "john doe",
			Phone:       "",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: phoneBook.ID,
			PhoneBook:   phoneBook,
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Phone is required", response["error"])
	})

	t.Run("CreatePhoneBookNumberInvalidPhone", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "john doe",
			Phone:       "1234",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: phoneBook.ID,
			PhoneBook:   phoneBook,
		}

		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Invalid Phone Number", response["error"])
	})

	t.Run("CreatePhoneBookNumberDuplicatePhone", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "john doe",
			Phone:       "09376304339",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: phoneBook.ID,
			PhoneBook:   phoneBook,
		}
		// in this test this object with this phone number is created in success test so should return error
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Input Phone Number has already been registered", response["error"])
	})

	t.Run("CreatePhoneBookNumberDuplicateUsername", func(t *testing.T) {
		requestBody := models.PhoneBookNumber{
			Name:        "john doe",
			Phone:       "09376304338",
			Prefix:      "1",
			Username:    "johndoe",
			PhoneBookID: phoneBook.ID,
			PhoneBook:   phoneBook,
		}
		// in this test this object with this username is created in success test so should return error
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.Set("account", account)

		err := handler.CreatePhoneBookNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var response map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Input Username has already been registered", response["error"])
	})
}

func TestListPhoneBookNumbers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
			Budget:   0,
			Password: "password",
			IsActive: true,
			IsAdmin:  false,
		}
		err = db.Create(&account).Error
		assert.NoError(t, err)
		phoneBook := models.PhoneBook{
			AccountID: account.ID,
			Name:      "Phone Book 1",
		}
		err = db.Create(&phoneBook).Error
		assert.NoError(t, err)

		phoneBookNumbers := []models.PhoneBookNumber{
			{
				Name:        "John Doe",
				Phone:       "123456789",
				Prefix:      "1",
				Username:    "johndoe",
				PhoneBookID: phoneBook.ID,
				PhoneBook:   phoneBook,
			},
			{
				Name:        "Jane Smith",
				Phone:       "987654321",
				Prefix:      "1",
				Username:    "janesmith",
				PhoneBookID: phoneBook.ID,
				PhoneBook:   phoneBook,
			},
		}
		err = db.Create(&phoneBookNumbers).Error
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/account/phone-books/"+fmt.Sprint(phoneBook.ID)+"/phone-book-numbers", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookID")
		ctx.SetParamValues(fmt.Sprint(phoneBook.ID))
		ctx.Set("account", account)
		handler := handlers.NewPhonebookHandler(db)
		err = handler.ListPhoneBookNumbers(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.PhoneBookNumber
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Len(t, response, len(phoneBookNumbers))
	})

	t.Run("PhonebookNotFound", func(t *testing.T) {
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
			Budget:   0,
			Password: "password",
			IsActive: true,
			IsAdmin:  false,
		}
		err = db.Create(&account).Error
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/account/phone-books/"+fmt.Sprint(9999)+"/phone-book-numbers", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookID")
		ctx.SetParamValues(fmt.Sprint(9999))
		ctx.Set("account", account)

		handler := handlers.NewPhonebookHandler(db)
		err = handler.ListPhoneBookNumbers(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	})
}

func TestReadPhoneBookNumber(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, err := utils.CreateTestDatabase()
		defer utils.CloseTestDatabase(db)
		assert.NoError(t, err)

		phoneBookNumber := models.PhoneBookNumber{
			Name:   "John Doe",
			Phone:  "123456789",
			Prefix: "1",
		}
		err = db.Create(&phoneBookNumber).Error
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/account/phone-books/phone-book-numbers/"+fmt.Sprint(phoneBookNumber.ID), nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookNumberID")
		ctx.SetParamValues(fmt.Sprint(phoneBookNumber.ID))

		handler := handlers.NewPhonebookHandler(db)
		err = handler.ReadPhoneBookNumber(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.PhoneBookNumber
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, phoneBookNumber.ID, response.ID)
		assert.Equal(t, phoneBookNumber.Name, response.Name)
		assert.Equal(t, phoneBookNumber.Phone, response.Phone)
		assert.Equal(t, phoneBookNumber.Prefix, response.Prefix)
	})

	t.Run("PhoneBookNumberNotFound", func(t *testing.T) {
		db, err := utils.CreateTestDatabase()
		defer utils.CloseTestDatabase(db)
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/account/phone-books/phone-book-numbers/"+fmt.Sprint(9999), nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookNumberID")
		ctx.SetParamValues(fmt.Sprint(9999))

		handler := handlers.NewPhonebookHandler(db)
		err = handler.ReadPhoneBookNumber(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Phone book number not found\"\n", rec.Body.String())
	})
}

func TestUpdatePhoneBookNumber(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, err := utils.CreateTestDatabase()
		defer utils.CloseTestDatabase(db)
		assert.NoError(t, err)

		existingPhoneBookNumber := models.PhoneBookNumber{
			Name:   "John Doe",
			Phone:  "123456789",
			Prefix: "1",
		}
		err = db.Create(&existingPhoneBookNumber).Error
		assert.NoError(t, err)

		e := echo.New()

		updatedPhoneBookNumber := handlers.UpdatePhoneBookNumberRequest{
			Prefix: "2",
			Name:   "Jane Smith",
			Phone:  "987654321",
		}
		jsonBody, _ := json.Marshal(updatedPhoneBookNumber)

		req := httptest.NewRequest(http.MethodPut, "/account/phone-books/phone-book-numbers/"+fmt.Sprint(existingPhoneBookNumber.ID), bytes.NewReader(jsonBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookNumberID")
		ctx.SetParamValues(fmt.Sprint(existingPhoneBookNumber.ID))

		handler := handlers.NewPhonebookHandler(db)
		err = handler.UpdatePhoneBookNumber(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.PhoneBookNumber
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, updatedPhoneBookNumber.Prefix, response.Prefix)
		assert.Equal(t, updatedPhoneBookNumber.Name, response.Name)
		assert.Equal(t, updatedPhoneBookNumber.Phone, response.Phone)

		// Verify that the object is updated in the database
		var updatedObject models.PhoneBookNumber
		err = db.First(&updatedObject, existingPhoneBookNumber.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, updatedPhoneBookNumber.Prefix, updatedObject.Prefix)
		assert.Equal(t, updatedPhoneBookNumber.Name, updatedObject.Name)
		assert.Equal(t, updatedPhoneBookNumber.Phone, updatedObject.Phone)
	})

	t.Run("PhoneBookNumberNotFound", func(t *testing.T) {
		db, err := utils.CreateTestDatabase()
		defer utils.CloseTestDatabase(db)
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodPut, "/account/phone-books/phone-book-numbers/"+fmt.Sprint(9999), nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookNumberID")
		ctx.SetParamValues(fmt.Sprint(9999))

		handler := handlers.NewPhonebookHandler(db)
		err = handler.UpdatePhoneBookNumber(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Phone book number not found\"\n", rec.Body.String())
	})
}

func TestDeletePhoneBookNumber(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, err := utils.CreateTestDatabase()
		defer utils.CloseTestDatabase(db)
		assert.NoError(t, err)

		existingPhoneBookNumber := models.PhoneBookNumber{
			Name:   "John Doe",
			Phone:  "123456789",
			Prefix: "1",
		}
		err = db.Create(&existingPhoneBookNumber).Error
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/account/phone-books/phone-book-numbers/"+fmt.Sprint(existingPhoneBookNumber.ID), nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookNumberID")
		ctx.SetParamValues(fmt.Sprint(existingPhoneBookNumber.ID))

		handler := handlers.NewPhonebookHandler(db)
		err = handler.DeletePhoneBookNumber(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"Phone book number deleted\"\n", rec.Body.String())

		var deletedPhoneBookNumber models.PhoneBookNumber
		result := db.First(&deletedPhoneBookNumber, existingPhoneBookNumber.ID)
		assert.Error(t, result.Error)
		assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
	})

	t.Run("PhoneBookNumberNotFound", func(t *testing.T) {
		db, err := utils.CreateTestDatabase()
		defer utils.CloseTestDatabase(db)
		assert.NoError(t, err)

		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/account/phone-books/phone-book-numbers/"+fmt.Sprint(9999), nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames("phoneBookNumberID")
		ctx.SetParamValues(fmt.Sprint(9999))

		handler := handlers.NewPhonebookHandler(db)
		err = handler.DeletePhoneBookNumber(ctx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Phone book number not found\"\n", rec.Body.String())
	})
}
