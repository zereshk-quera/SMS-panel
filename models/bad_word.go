package models

type Bad_Word struct {
	ID    uint   `gorm:"primaryKey"`
	Word  string `gorm:"unique"`
	Regex string
}
