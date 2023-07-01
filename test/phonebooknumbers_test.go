package test

/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestCreatePhoneBookNumber(t *testing.T) {
	tests := []struct {
		name           string
		phoneNumberReq models.PhoneBookNumber
		expectedCode   int
		expectedError  string
	}{
		{
			name: "Success",
			phoneNumberReq: models.PhoneBookNumber{
				PhoneBookID: phoneBookID,
				Username:    "john",
				Prefix:      "+98",
				Name:        "John Doe",
				Phone:       "09123456789",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "MissingName",
			phoneNumberReq: models.PhoneBookNumber{
				PhoneBookID: 1,
				Username:    "john2",
				Prefix:      "+1",
				// Name:        "John Doe",
				Phone: "9123456781",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Name is required",
		},
		{
			name: "MissingPhone",
			phoneNumberReq: models.PhoneBookNumber{
				PhoneBookID: 1,
				Username:    "john1",
				Prefix:      "+1",
				Name:        "John Doe",
				// Phone:       "23456789",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Phone is required",
		},
		{
			name: "InvalidPhone",
			phoneNumberReq: models.PhoneBookNumber{
				PhoneBookID: 1,
				Username:    "john1",
				Prefix:      "+1",
				Name:        "John Doe",
				Phone:       "23456789",
			},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Invalid Phone Number",
		},
		{
			name: "DuplicatePhone",
			phoneNumberReq: models.PhoneBookNumber{
				// PhoneBookID: 1,
				Username: "john1",
				Prefix:   "+98",
				Name:     "John Doe",
				Phone:    "09123456789",
			},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Inupt Phone Number has already been registered",
		},
		{
			name: "DuplicateUsername",
			phoneNumberReq: models.PhoneBookNumber{
				// PhoneBookID: 1,
				Username: "john",
				Prefix:   "+1",
				Name:     "John Doe",
				Phone:    "09123456782",
			},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Inupt Username has already been registered",
		},
	}

	for i, test := range tests {
		t.Skip("skipping for now")
		t.Run(test.name, func(t *testing.T) {
			reqBody, err := json.Marshal(test.phoneNumberReq)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/account/phone-books/phone-book-numbers", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = phonebookHandler.CreatePhoneBookNumber(c)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)

			if test.expectedCode == http.StatusCreated {
				var phoneNumberRes models.PhoneBookNumber
				err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
				assert.NoError(t, err)

				expectedResp := test.phoneNumberReq
				expectedResp.ID = phoneNumberRes.ID

				phoneBookNumberID = phoneNumberRes.ID

				if !reflect.DeepEqual(phoneNumberRes, expectedResp) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, expectedResp)
				}
			} else {
				var errorRes map[string]string
				err = json.Unmarshal(rec.Body.Bytes(), &errorRes)
				assert.NoError(t, err)

				if !reflect.DeepEqual(errorRes["error"], test.expectedError) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, errorRes["error"], test.expectedError)
				}
			}
		})
	}
}

func TestListPhoneBookNumbers(t *testing.T) {
	tests := []struct {
		name           string
		phoneNumberReq []models.PhoneBookNumber
		expectedCode   int
		expectedError  string
	}{
		{
			name: "Success",
			phoneNumberReq: []models.PhoneBookNumber{
				{
					PhoneBookID: phoneBookID,
					Username:    "john",
					Prefix:      "+98",
					Name:        "John Doe",
					Phone:       "09123456789",
				},
			},
			expectedCode: http.StatusOK,
		},
	}

	for i, test := range tests {
		t.Skip("skipping for now")
		t.Run(test.name, func(t *testing.T) {
			target := fmt.Sprintf("/account/phone-books/%d/phone-book-numbers", phoneBookID)
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("phoneBookID")
			c.SetParamValues(fmt.Sprint(phoneBookID))

			err := phonebookHandler.ListPhoneBookNumbers(c)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)

			if test.expectedCode == http.StatusOK {
				var phoneNumberRes []models.PhoneBookNumber
				err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
				assert.NoError(t, err)
				expectedResp := test.phoneNumberReq

				for i2, pbn := range phoneNumberRes {
					expectedResp[i2].ID = pbn.ID
				}

				if !reflect.DeepEqual(phoneNumberRes, expectedResp) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, expectedResp)
				}
			}
		})
	}
}

func TestReadPhoneBookNumbers(t *testing.T) {
	tests := []struct {
		name           string
		phoneNumberReq models.PhoneBookNumber
		expectedCode   int
		expectedError  string
	}{
		{
			name: "Success",
			phoneNumberReq: models.PhoneBookNumber{
				ID:          phoneBookNumberID,
				PhoneBookID: phoneBookID,
				Username:    "john",
				Prefix:      "+98",
				Name:        "John Doe",
				Phone:       "09123456789",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "NotFound",
			phoneNumberReq: models.PhoneBookNumber{
				ID: 500,
				// PhoneBookID: phoneBookID,
				// Username:    "john",
				// Prefix:      "+98",
				// Name:        "John Doe",
				// Phone:       "09123456789",
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "Phone book number not found",
		},
	}

	for i, test := range tests {
		t.Skip("skipping for now")
		t.Run(test.name, func(t *testing.T) {
			target := fmt.Sprintf("/account/phone-books/phone-book-numbers/%d", test.phoneNumberReq.ID)
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("phoneBookNumberID")
			c.SetParamValues(fmt.Sprint(test.phoneNumberReq.ID))

			err := phonebookHandler.ReadPhoneBookNumber(c)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)

			if test.expectedCode == http.StatusOK {
				var phoneNumberRes models.PhoneBookNumber
				err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
				assert.NoError(t, err)

				if !reflect.DeepEqual(phoneNumberRes, test.phoneNumberReq) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, test.phoneNumberReq)
				}
			} else {
				var phoneNumberRes string
				err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
				assert.NoError(t, err)

				if !reflect.DeepEqual(phoneNumberRes, test.expectedError) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, test.expectedError)
				}
			}
		})
	}
}

func TestUpdatePhoneBookNumbers(t *testing.T) {
	tests := []struct {
		name           string
		phoneNumberReq models.PhoneBookNumber
		expectedCode   int
		expectedError  string
	}{
		{
			name: "Success",
			phoneNumberReq: models.PhoneBookNumber{
				ID:          phoneBookNumberID,
				PhoneBookID: phoneBookID,
				Username:    "john",
				Prefix:      "+98",
				Name:        "John Doe 2",
				Phone:       "09123456789",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "NotFound",
			phoneNumberReq: models.PhoneBookNumber{
				ID:          500,
				PhoneBookID: phoneBookID,
				Username:    "john",
				Prefix:      "+98",
				Name:        "John Doe",
				Phone:       "09123456789",
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "Phone book number not found",
		},
	}

	for i, test := range tests {
		t.Skip("skipping for now")
		t.Run(test.name, func(t *testing.T) {
			reqBody, err := json.Marshal(test.phoneNumberReq)
			assert.NoError(t, err)

			target := fmt.Sprintf("/account/phone-books/phone-book-numbers/%d", test.phoneNumberReq.ID)
			req := httptest.NewRequest(http.MethodPut, target, bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("phoneBookNumberID")
			c.SetParamValues(fmt.Sprint(test.phoneNumberReq.ID))

			err = phonebookHandler.UpdatePhoneBookNumber(c)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)

			if test.expectedCode == http.StatusOK {
				var phoneNumberRes models.PhoneBookNumber
				err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
				assert.NoError(t, err)

				if !reflect.DeepEqual(phoneNumberRes, test.phoneNumberReq) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, test.phoneNumberReq)
				}
			} else {
				var phoneNumberRes string
				err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
				assert.NoError(t, err)

				if !reflect.DeepEqual(phoneNumberRes, test.expectedError) {
					t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, test.expectedError)
				}
			}
		})
	}
}

func TestDeletePhoneBookNumbers(t *testing.T) {
	tests := []struct {
		name           string
		phoneNumberReq models.PhoneBookNumber
		expectedCode   int
		expectedResp   string
	}{
		{
			name: "Success",
			phoneNumberReq: models.PhoneBookNumber{
				ID: phoneBookNumberID,
			},
			expectedCode: http.StatusOK,
			expectedResp: "Phone book number deleted",
		},
		{
			name: "NotFound",
			phoneNumberReq: models.PhoneBookNumber{
				ID: 500,
			},
			expectedCode: http.StatusNotFound,
			expectedResp: "Phone book number not found",
		},
	}

	for i, test := range tests {
		t.Skip("skipping for now")
		t.Run(test.name, func(t *testing.T) {
			reqBody, err := json.Marshal(test.phoneNumberReq)
			assert.NoError(t, err)

			target := fmt.Sprintf("/account/phone-books/phone-book-numbers/%d", test.phoneNumberReq.ID)
			req := httptest.NewRequest(http.MethodDelete, target, bytes.NewReader(reqBody))
			// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("phoneBookNumberID")
			c.SetParamValues(fmt.Sprint(test.phoneNumberReq.ID))

			err = phonebookHandler.DeletePhoneBookNumber(c)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedCode, rec.Code)

			var phoneNumberRes string
			err = json.Unmarshal(rec.Body.Bytes(), &phoneNumberRes)
			assert.NoError(t, err)

			if !reflect.DeepEqual(phoneNumberRes, test.expectedResp) {
				t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, phoneNumberRes, test.expectedResp)
			}
		})
	}
}
*/
