package model

import "gorm.io/gorm"

type Pet struct {
	gorm.Model
	Phone                 string
	PetName               string
	PetType               string
	Breed                 string
	PetIns                string
	Birthday              string
	Gender                string
	AdditionalInformation string
	ShopifyCustomerID     string `gorm:"size:64;not null;index"`
}

// 创建表
func MigratePet(db *gorm.DB) error {
	return db.AutoMigrate(&Pet{})
}
