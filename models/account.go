package models

type Account struct {
	ID       uint   `gorm:"primary_key"`
	UserID   uint   `gorm:"not null"`
	Username string `gorm:"type:varchar(255);unique;not null"`
	Budget   int64  `gorm:"type:bigint"`
	Password string `gorm:"type:varchar(255)"`
	Token    string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
}
