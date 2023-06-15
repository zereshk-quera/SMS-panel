package test

import (
	database "SMS-panel/database"
	"SMS-panel/handlers"
	"SMS-panel/models"

)

var (
	e *echo.Echo
)

// @Router /accounts/register
func TestRegisterHandler(t *testing.T) {
	t.Run("Stupid User without FirstName input", func(t *testing.T) {
		
		withoutFirstNameUserCreate := handlers.UserCreateRequest{
			//FirstName: "Rick",
			LastName: "Sanchez",  		
			Email: "RickSanchez@morty.com",
			Phone: "09123456789",		
			NationalID: "0369734971",
			Username: "ricksanchez",
			Password: "123Rick123"  		
		}

		reqBody, err := json.Marshal(withoutFirstNameUserCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/accounts/register",bytes.NewReader(reqBody))
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
			Email: "RickSanchez@morty.com",
			Phone: "09123456789",		
			NationalID: "0369734971",
			Username: "ricksanchez",
			Password: "123Rick123"  		
		}

		reqBody, err := json.Marshal(withoutLastNameUserCreate)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/accounts/register",bytes.NewReader(reqBody))
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
}
