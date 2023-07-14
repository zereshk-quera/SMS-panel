package models

type Budget struct {
	ID        uint  `gorm:"primary_key"`
	AccountID uint  `gorm:"not null"`
	Amount    int64 `gorm:"type:bigint"`
}
