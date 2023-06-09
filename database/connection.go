package db

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Connect() error {
	dsn := "host=your_host port=your_port user=your_user password=your_password dbname=your_database sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	dbConn = db
	return nil
}

func GetConnection() (*gorm.DB, error) {
	if dbConn == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return dbConn, nil
}
