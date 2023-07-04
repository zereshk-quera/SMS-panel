package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	database "SMS-panel/database"
	"SMS-panel/middlewares"
	"SMS-panel/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	deleteTestAccount()
	os.Exit(exitCode)
}

func TestIsLoggedIn(t *testing.T) {
	e := echo.New()
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	}, middlewares.IsLoggedIn)

	token, err := addTestAccount()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", token)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	err = deleteTestAccount()
	assert.NoError(t, err)
}

func addTestAccount() (string, error) {
	db, err := database.GetConnection()
	if err != nil {
		return "", err
	}

	user := models.User{
		FirstName:  "john",
		LastName:   "doe",
		Phone:      "09376304339",
		Email:      "test@gmail.com",
		NationalID: "123456789",
	}
	err = db.Create(&user).Error
	if err != nil {
		return "", err
	}

	account := models.Account{
		UserID:   user.ID,
		Username: "testuser",
		Budget:   10,
		Password: "password",
		IsActive: true,
		IsAdmin:  false,
	}

	err = db.Create(&account).Error
	if err != nil {
		return "", err
	}

	token, err := createTokenForAccount(account)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createTokenForAccount(account models.Account) (string, error) {
	claims := jwt.MapClaims{
		"id":  account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func deleteTestAccount() error {
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	err = db.Delete(models.Account{}, "username = ?", "testuser").Error
	if err != nil {
		return err
	}

	return nil
}
