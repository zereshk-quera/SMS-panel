package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Connect() error {
	dsn := "host=localhost port=5433 user=amirhejazi password=postgres dbname=sms_panel sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
			return nil, err
		}
	}
	return dbConn, nil
}
