package models

type PhoneBook struct {
	ID        uint   `gorm:"primary_key"`
	AccountID uint   `gorm:"not null"`
	Name      string `gorm:"type:varchar(255)"`
}
type PhoneBookNumber struct {
	ID          uint   `gorm:"primary_key"`
	PhoneBookID uint   `gorm:"not null"`
	Prefix      string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	Phone       string `gorm:"type:varchar(255)"`
}
