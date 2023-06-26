package db

import (
	"errors"
	"fmt"
	"os"
	"time"

	"SMS-panel/models"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbConn *gorm.DB

func Connect() error {
	// -------env----------
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
	// 	cfg.PG.HOST, cfg.PG.USER, cfg.PG.PASSWORD, cfg.PG.DB, cfg.PG.PORT, cfg.PG.SSLMODE, cfg.PG.TIMEZONE)
	// -------env----------

	// If not connect - use "db" instead of "localhost"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		"localhost", "amirhejazi", "postgres", "sms_panel", "5433", "disable", "Asia/Tehran")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	dbConn = db
	return nil
}

func GetConnection() (*gorm.DB, error) {
	if dbConn == nil {
		err := Connect()
		if err != nil {
			return nil, errors.New("database connection is not initialized")
		}
	}
	// check super admin existence
	var account models.Account
	dbConn.Where("id = ?", 1).First(&account)

	// initial super admin
	if account.ID == 0 {
		var user models.User
		user.FirstName = "admin"
		user.LastName = "admin"
		user.Phone = "admin"
		user.Email = "admin"
		user.NationalID = "admin"
		dbConn.Save(&user)
		var adminAccount models.Account
		adminAccount.UserID = 1
		adminAccount.Budget = 0
		adminAccount.IsActive = true
		adminAccount.IsAdmin = true
		adminAccount.Username = "admin"
		// save hashed password
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		adminAccount.Password = string(hash)
		// Generate Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  1,
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
		adminAccount.Token = tokenString
		dbConn.Save(&adminAccount)
	}

	return dbConn, nil
}
