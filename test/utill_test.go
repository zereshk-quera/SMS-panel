package test

import (
	"testing"

	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateEmail(t *testing.T) {
	t.Run("ValidEmail", func(t *testing.T) {
		validEmail := "test@example.com"
		isValid := utils.ValidateEmail(validEmail)
		assert.True(t, isValid, "Expected valid email")
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		invalidEmail := "invalid-email"
		isValid := utils.ValidateEmail(invalidEmail)
		assert.False(t, isValid, "Expected invalid email")
	})
}

func TestValidatePhone(t *testing.T) {
	t.Run("ValidPhone", func(t *testing.T) {
		validPhone := "09376304339"
		isValid := utils.ValidatePhone(validPhone)
		assert.True(t, isValid, "Expected valid phone")
	})

	t.Run("InvalidPhoneWithCharacter", func(t *testing.T) {
		invalidPhone := "0937abc4339"
		isValid := utils.ValidatePhone(invalidPhone)
		assert.False(t, isValid, "Expected invalid phone with character")
	})

	t.Run("InvalidPhoneLength", func(t *testing.T) {
		invalidPhone := "0937630433"
		isValid := utils.ValidatePhone(invalidPhone)
		assert.False(t, isValid, "Expected invalid phone length")
	})

	t.Run("InvalidPhonePrefix", func(t *testing.T) {
		invalidPhone := "08376304339"
		isValid := utils.ValidatePhone(invalidPhone)
		assert.False(t, isValid, "Expected invalid phone prefix")
	})
}

func TestParseInt(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		validInput := "123"
		expectedResult := 123
		result := utils.ParseInt(validInput, 10)
		assert.Equal(t, expectedResult, result, "Expected valid input to be parsed correctly")
	})

	t.Run("InvalidInput", func(t *testing.T) {
		invalidInput := "abc"
		expectedResult := 0
		result := utils.ParseInt(invalidInput, 10)
		assert.Equal(t, expectedResult, result, "Expected invalid input to return 0")
	})
}

func TestValidateNationalID(t *testing.T) {
	t.Run("ValidNationalID", func(t *testing.T) {
		validNationalID := "0817762590"
		result := utils.ValidateNationalID(validNationalID)
		assert.True(t, result, "Expected valid national ID to be validated as true")
	})

	t.Run("InvalidNationalID", func(t *testing.T) {
		invalidNationalID := "1234567890"
		result := utils.ValidateNationalID(invalidNationalID)
		assert.False(t, result, "Expected invalid national ID to be validated as false")
	})
}

func TestValidateTimeFormat(t *testing.T) {
	t.Run("ValidTimeFormat", func(t *testing.T) {
		validTime := "12:34"
		result := utils.ValidateTimeFormat(validTime)
		assert.True(t, result, "Expected valid time format to be validated as true")
	})

	t.Run("InvalidTimeFormat", func(t *testing.T) {
		invalidTime := "12:345"
		result := utils.ValidateTimeFormat(invalidTime)
		assert.False(t, result, "Expected invalid time format to be validated as false")
	})
}

func TestValidateJsonFormat(t *testing.T) {
	t.Run("ValidJsonFormat", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}
		fields := []string{"field1", "field2"}
		result, err := utils.ValidateJsonFormat(jsonBody, fields...)
		assert.NoError(t, err, "Expected valid JSON format to have no error")
		assert.Equal(t, "OK", result, "Expected valid JSON format to return 'OK'")
	})

	t.Run("InvalidJsonFormat", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"field1": "value1",
		}
		fields := []string{"field1", "field2"}
		result, err := utils.ValidateJsonFormat(jsonBody, fields...)
		assert.Error(t, err, "Expected invalid JSON format to have an error")
		assert.Contains(t, result, "Input Json doesn't include field2", "Expected invalid JSON format to return error message")
	})
}

func TestValidateUser(t *testing.T) {
	t.Run("ValidUser", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "0817762590",
		}

		msg, user, err := utils.ValidateUser(jsonBody)

		assert.NoError(t, err)
		assert.Equal(t, "OK", msg)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "johndoe@example.com", user.Email)
		assert.Equal(t, "09376304339", user.Phone)
		assert.Equal(t, "0817762590", user.NationalID)
	})

	t.Run("InvalidFirstName", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"firstname":  "",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "0817762590",
		}

		msg, _, err := utils.ValidateUser(jsonBody)

		assert.EqualError(t, err, "")
		assert.Equal(t, "First Name can't be empty", msg)
	})

	t.Run("InvalidLastName", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "0817762590",
		}

		msg, _, err := utils.ValidateUser(jsonBody)

		assert.EqualError(t, err, "")
		assert.Equal(t, "Last Name can't be empty", msg)
	})

	t.Run("InvalidPhoneNumber", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "12345",
			"nationalid": "0817762590",
		}

		msg, _, err := utils.ValidateUser(jsonBody)

		assert.EqualError(t, err, "")
		assert.Equal(t, "Invalid Phone Number", msg)
	})

	t.Run("InvalidEmailAddress", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoecom",
			"phone":      "09376304339",
			"nationalid": "0817762590",
		}

		msg, _, err := utils.ValidateUser(jsonBody)

		assert.EqualError(t, err, "")
		assert.Equal(t, "Invalid Email Address", msg)
	})

	t.Run("InvalidNationalID", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"firstname":  "John",
			"lastname":   "Doe",
			"email":      "johndoe@example.com",
			"phone":      "09376304339",
			"nationalid": "1234567890",
		}

		msg, _, err := utils.ValidateUser(jsonBody)

		assert.EqualError(t, err, "")
		assert.Equal(t, "Invalid National ID", msg)
	})
}

func TestCheckUnique(t *testing.T) {
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	t.Run("UniqueData", func(t *testing.T) {
		user := models.User{
			Phone:      "123456789",
			Email:      "test1@example.com",
			NationalID: "1234567890",
		}
		username := "uniqueusername"

		msg, err := utils.CheckUnique(user, username, db)

		assert.NoError(t, err)
		assert.Equal(t, "OK", msg)
	})

	t.Run("NonUniquePhone", func(t *testing.T) {
		existingUser := models.User{
			Phone:      "1234567890",
			Email:      "test2@example.com",
			NationalID: "12345678900",
		}
		db.Create(&existingUser)

		user := models.User{
			Phone: "1234567890",
		}
		username := "uniqueusername"

		msg, err := utils.CheckUnique(user, username, db)

		assert.Error(t, err)
		assert.Equal(t, "Input Phone Number has already been registered", msg)
	})

	t.Run("NonUniqueEmail", func(t *testing.T) {
		existingUser := models.User{
			Phone:      "12345678900",
			Email:      "test@example.com",
			NationalID: "123456789000",
		}
		db.Create(&existingUser)

		user := models.User{
			Email: "test@example.com",
		}
		username := "uniqueusername"

		msg, err := utils.CheckUnique(user, username, db)

		assert.Error(t, err)
		assert.Equal(t, "Input Email Address has already been registered", msg)
	})

	t.Run("NonUniqueNationalID", func(t *testing.T) {
		existingUser := models.User{
			NationalID: "123456789",
			Email:      "test3@example.com",
			Phone:      "1234567190000",
		}
		db.Create(&existingUser)

		user := models.User{
			NationalID: "123456789",
		}
		username := "uniqueusername"

		msg, err := utils.CheckUnique(user, username, db)

		assert.Error(t, err)
		assert.Equal(t, "Input National ID has already been registered", msg)
	})

	t.Run("NonUniqueUsername", func(t *testing.T) {
		existingAccount := models.Account{
			Username: "uniqueusername",
		}
		db.Create(&existingAccount)

		user := models.User{}
		username := "uniqueusername"

		msg, err := utils.CheckUnique(user, username, db)

		assert.Error(t, err)
		assert.Equal(t, "Input Username has already been registered", msg)
	})
}

func TestGenerateRegex(t *testing.T) {
	input := "hi"
	expected := "[Hh#ʜɦλ][^a-zA-Z]*[Ii1!|ɨ]"

	result := utils.GenerateRegex(input)
	if result != expected {
		t.Errorf("Input: %s\nExpected: %s\nGot: %s", input, expected, result)
	}
}

func TestCreateAccount(t *testing.T) {
	// Create a test database
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err, "Failed to create test database")
	defer func() {
		err := utils.CloseTestDatabase(db)
		assert.NoError(t, err, "Failed to close test database")
	}()

	t.Run("SuccessfulAccountCreation", func(t *testing.T) {
		userID := 123
		username := "testuser"
		isAdmin := false
		password := "password"

		msg, account, err := utils.CreateAccount(userID, username, isAdmin, password, db)

		assert.NoError(t, err, "Expected no error in creating account")
		assert.Equal(t, "OK", msg, "Expected success message to be 'OK'")
		assert.NotZero(t, account.ID, "Expected account ID to be non-zero")
		assert.Equal(t, uint(userID), account.UserID, "Expected account UserID to match the provided value")
		assert.Equal(t, username, account.Username, "Expected account Username to match the provided value")
		assert.Equal(t, int64(0), account.Budget, "Expected account Budget to be zero")
		assert.Equal(t, false, account.IsAdmin, "Expected account IsAdmin to match the provided value")
		assert.NotEmpty(t, account.Token, "Expected account Token to be non-empty")
	})

	t.Run("FailedAccountCreation", func(t *testing.T) {
		// Provide invalid test data
		userID := 123
		username := "testuser"
		isAdmin := false
		password := ""

		msg, account, err := utils.CreateAccount(userID, username, isAdmin, password, db)

		assert.Error(t, err, "Expected error in creating account")
		assert.Equal(t, "Failed to Create Account", msg, "Failed to Create Account")
		assert.Empty(t, account, "Expected empty account")
	})
}

func TestLogin(t *testing.T) {
	// Create a test database
	db, err := utils.CreateTestDatabase()
	assert.NoError(t, err)
	defer utils.CloseTestDatabase(db)

	user := models.User{FirstName: "testuser", LastName: "testuser", Phone: "09376304339", Email: "amir@gmail.com", NationalID: "0265670578"}
	db.Create(&user)
	hash, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	assert.NoError(t, err)
	account := models.Account{UserID: user.ID, Username: "testuser", Password: string(hash), Token: "testtoken", IsActive: true, IsAdmin: false}
	db.Create(&account)

	t.Run("ValidLogin", func(t *testing.T) {
		msg, resultAccount, err := utils.Login("testuser", "test123", false, db)

		assert.NoError(t, err, "Expected no error in Login")
		assert.Equal(t, "OK", msg, "Expected message to be 'OK'")
		assert.Equal(t, account.ID, resultAccount.ID, "Expected account ID to match")
		assert.NotEmpty(t, resultAccount.Token, "Expected account token to be generated")
	})

	t.Run("InvalidUsername", func(t *testing.T) {
		msg, _, err := utils.Login("invaliduser", "test123", false, db)

		assert.Error(t, err, "Expected error in Login")
		assert.Equal(t, "Invalid Username", msg, "Expected message to be 'Invalid Username'")
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		msg, _, err := utils.Login("testuser", "wrongpassword", false, db)

		assert.Error(t, err, "Expected error in Login")
		assert.Equal(t, "Wrong Password", msg, "Expected message to be 'Wrong Password'")
	})

	t.Run("InactiveAccount", func(t *testing.T) {
		account.IsActive = false
		db.Save(&account)

		msg, _, err := utils.Login("testuser", "test123", false, db)

		assert.Error(t, err, "Expected error in Login")
		assert.Equal(t, "Your Account Isn't Active", msg, "Expected message to be 'Your Account Isn't Active'")
	})

	t.Run("AdminLogin", func(t *testing.T) {
		account.IsAdmin = true
		account.IsActive = true
		db.Save(&account)

		msg, resultAccount, err := utils.Login("testuser", "test123", true, db)

		assert.NoError(t, err, "Expected no error in Login")
		assert.Equal(t, "OK", msg, "Expected message to be 'OK'")
		assert.Equal(t, account.ID, resultAccount.ID, "Expected account ID to match")
		assert.NotEmpty(t, resultAccount.Token, "Expected account token to be generated")
		assert.True(t, resultAccount.IsAdmin, "Expected account to be admin")
	})
}
