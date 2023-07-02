package utils

import (
	"SMS-panel/models"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// this function used to check user properties validation
func ValidateUser(jsonBody map[string]interface{}) (string, models.User, error) {
	msg := "OK"
	// Create User Object
	var user models.User
	user.FirstName = jsonBody["firstname"].(string)
	user.LastName = jsonBody["lastname"].(string)
	user.Email = jsonBody["email"].(string)
	user.Phone = jsonBody["phone"].(string)
	user.NationalID = jsonBody["nationalid"].(string)

	// Check FirstName Validation
	if len(strings.TrimSpace(user.FirstName)) == 0 {
		msg = "First Name can't be empty"
		return msg, models.User{}, errors.New("")
	}

	// Check LastName Validation
	if len(strings.TrimSpace(user.LastName)) == 0 {
		msg = "Last Name can't be empty"
		return msg, models.User{}, errors.New("")
	}

	// Check Phone Number Validation
	if !ValidatePhone(user.Phone) {
		msg = "Invalid Phone Number"
		return msg, models.User{}, errors.New("")
	}

	// Check Email Validation
	if !ValidateEmail(user.Email) {
		msg = "Invalid Email Address"
		return msg, models.User{}, errors.New("")
	}

	// Check NationalID Validation
	if !ValidateNationalID(user.NationalID) {
		msg = "Invalid National ID"
		return msg, models.User{}, errors.New("")
	}
	return msg, user, nil
}

func CheckUnique(user models.User, username string, db *gorm.DB) (string, error) {
	msg := "OK"
	// Is Input Phone Number Unique or Not
	var existingUser models.User
	db.Where("phone = ?", user.Phone).First(&existingUser)
	if existingUser.ID != 0 {
		msg = "Inupt Phone Number has already been registered"
		return msg, errors.New("")
	}

	// Is Input Email Address Unique or Not
	db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		msg = "Inupt Email Address has already been registered"
		return msg, errors.New("")
	}

	// Is Input National ID Unique or Not
	db.Where("national_id = ?", user.NationalID).First(&existingUser)
	if existingUser.ID != 0 {
		msg = "Inupt National ID has already been registered"
		return msg, errors.New("")
	}

	// Is Input Username Unique or Not
	var existingAccount models.Account
	db.Where("username = ?", username).First(&existingAccount)
	if existingAccount.ID != 0 {
		msg = "Inupt Username has already been registered"
		return msg, errors.New("")
	}
	return msg, nil
}

// this function used to create an account and insert it into database
func CreateAccount(user_id int, username string, is_admin bool, password string, db *gorm.DB) (string, models.Account, error) {
	msg := "OK"
	// Instantiating Account Object
	var account models.Account
	account.UserID = uint(user_id)
	account.Username = username
	account.Budget = 0
	account.IsAdmin = is_admin
	account.Token = ""

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		msg = "Failed to Hashing Password"
		return msg, models.Account{}, errors.New("")
	}
	account.Password = string(hash)

	//insert account into database
	createdAccount := db.Create(&account)
	if createdAccount.Error != nil {
		msg = "Failed to Create Account"
		return msg, models.Account{}, errors.New("")
	}

	//generate token
	var token *jwt.Token
	if is_admin {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    account.ID,
			"exp":   time.Now().Add(time.Hour).Unix(),
			"admin": true,
		})
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    account.ID,
			"exp":   time.Now().Add(time.Hour).Unix(),
			"admin": false,
		})
	}
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		msg = "Failed To Create Token"
		return msg, models.Account{}, errors.New("")
	}
	account.Token = tokenString

	//update account
	db.Save(&account)

	return msg, account, nil

}

func Login(username, password string, is_admin bool, db *gorm.DB) (string, models.Account, error) {
	msg := "OK"
	// Find account based on input username
	var account models.Account

	db.Where("username = ?", username).First(&account)

	// Account Not Found
	if account.ID == 0 {
		msg = "Invalid Username"
		return msg, models.Account{}, errors.New("")
	}

	// Incorrect Password
	err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		msg = "Wrong Password"
		return msg, models.Account{}, errors.New("")
	}

	if is_admin {
		if !account.IsAdmin {
			msg = "You are not admin!"
			return msg, models.Account{}, errors.New("")
		}
	}
	//Account isn't active
	if !account.IsActive {
		msg = "Your Account Isn't Active"
		return msg, models.Account{}, errors.New("")
	}

	var token *jwt.Token
	if is_admin {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    account.ID,
			"exp":   time.Now().Add(time.Hour).Unix(),
			"admin": true,
		})
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":    account.ID,
			"exp":   time.Now().Add(time.Hour).Unix(),
			"admin": false,
		})
	}

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		msg = "Failed To Create Token"
		return msg, models.Account{}, errors.New("")
	}

	// Update Account's Token In Database
	account.Token = tokenString
	db.Save(&account)
	return msg, account, nil
}
