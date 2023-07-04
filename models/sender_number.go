package models

type SenderNumber struct {
	ID         uint   `gorm:"primary_key"`
	Number     string `gorm:"type:varchar(255); unique; not null"`
	IsExlusive bool   `gorm:"type:bool; default:false; not null"`
	IsDefault  bool   `gorm:"type:bool; default:false; not null"`
}

func (SenderNumber) TableName() string {
	return "sender_numbers"
}
