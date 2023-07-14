package models

type PhoneBook struct {
	ID        uint   `gorm:"primary_key"`
	AccountID uint   `gorm:"not null"`
	Name      string `gorm:"type:varchar(255)"`
}

type PhoneBookNumber struct {
	ID          uint      `gorm:"primary_key"`
	PhoneBookID uint      `gorm:"not null"`
	Username    string    `gorm:"type:varchar(255);unique;default:null"`
	Prefix      string    `gorm:"type:varchar(255);default:+98"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Phone       string    `gorm:"type:varchar(255);unique;not null"`
	PhoneBook   PhoneBook `gorm:"association_autoupdate:false"`
}
