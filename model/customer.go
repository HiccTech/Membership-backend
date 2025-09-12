package model

import (
	"gorm.io/gorm"
)

type Customer struct {
	BaseModelNoID
	ShopifyCustomerId string  `gorm:"primaryKey;size:64"`
	Pets              []Pet   `gorm:"foreignKey:ShopifyCustomerId;references:ShopifyCustomerId"`
	Topups            []Topup `gorm:"foreignKey:ShopifyCustomerId;references:ShopifyCustomerId"`
}

func MigrateCustomer(db *gorm.DB) error {
	return db.AutoMigrate(&Customer{})
}
