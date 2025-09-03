package model

import (
	"gorm.io/gorm"
)

type Customer struct {
	BaseModelNoID
	ShopifyCustomerID string `gorm:"primaryKey;size:64"`
	Pets              []Pet  `gorm:"foreignKey:ShopifyCustomerID;references:ShopifyCustomerID"`
}

func MigrateCustomer(db *gorm.DB) error {
	return db.AutoMigrate(&Customer{})
}
