package test

import (
	"os"
	"testing"

	database "SMS-panel/database"
	"SMS-panel/handlers"
	"SMS-panel/models"

	"github.com/labstack/echo/v4"
)

var (
	account           models.Account
	phoneBookID       uint
	phoneBookNumberID uint
	e                 *echo.Echo
	phonebookHandler  *handlers.PhonebookHandler
)

func TestMain(m *testing.M) {
	err := database.Connect()
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

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
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	user := models.User{
		FirstName:  "John",
		LastName:   "Doe",
		Phone:      "123456789",
		Email:      "john.doe@example.com",
		NationalID: "1234567890",
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

	err = db.Exec("DELETE FROM phone_book_numbers WHERE phone_book_id = ?", phoneBookID).Error
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
