package model

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ShopifyCustomerID string `gorm:"primaryKey;size:64"`
	Pets              []Pet  `gorm:"foreignKey:ShopifyCustomerID;references:ShopifyCustomerID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func MigrateCustomer(db *gorm.DB) error {
	return db.AutoMigrate(&Customer{})
}
