package utils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"SMS-panel/models"
)

func CreateTestDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Account{}, &models.Bad_Word{},
		&models.Configuration{}, &models.PhoneBook{}, &models.PhoneBookNumber{},
		&models.Transaction{}, &models.SMSMessage{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseTestDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
