package db

import (
	"SMS-panel/config"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := cfg.PG.DSN
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