package models

import "time"

type UserNumbers struct {
	ID          uint      `gorm:"primary_key"`
	UserID      uint      `gorm:"not null"`
	NumberID    uint      `gorm:"not null"`
	StartDate   time.Time `gorm:"type:date"`
	EndDate     time.Time `gorm:"type:date"`
	IsAvailable bool      `gorm:"default:true"`
	User        User
	Number      SenderNumber `gorm:"foreignKey:NumberID"`
}

func (UserNumbers) TableName() string {
	return "user_numbers"
}
