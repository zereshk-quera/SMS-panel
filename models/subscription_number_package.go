package models

type SubscriptionNumberPackage struct {
	ID    uint   `gorm:"primary_key"`
	Title string `gorm:"type:varchar(55);not null;unique"`
	Price int64  `gorm:"type:bigint"`
}

func (SubscriptionNumberPackage) TableName() string {
	return "subscription_number_package"
}
