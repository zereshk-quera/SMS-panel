package models

import (
	"time"
)

type SMSMessage struct {
	ID             uint       `gorm:"primary_key"`
	Sender         string     `gorm:"type:varchar(255);not null"`
	Recipient      string     `gorm:"type:varchar(255);not null"`
	Message        string     `gorm:"type:text;not null"`
	Schedule       *time.Time `gorm:"default:null"`
	DeliveryReport string     `gorm:"type:text"`
	CreatedAt      time.Time  `gorm:"default:current_timestamp"`
	AccountID      uint       `gorm:"not null"`
}
