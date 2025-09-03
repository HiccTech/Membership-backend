package model

import "gorm.io/gorm"

type Pet struct {
	BaseModel
	Phone                 string `json:"phone"`
	PetName               string `json:"petName"`
	PetType               string `json:"petType"`
	Breed                 string `json:"breed"`
	PetIns                string `json:"petIns"`
	Birthday              string `json:"birthday"`
	Gender                string `json:"gender"`
	AdditionalInformation string `json:"additionalInformation"`
	ShopifyCustomerID     string `gorm:"size:64;not null;index" json:"shopifyCustomerId"`
}

// 创建表
func MigratePet(db *gorm.DB) error {
	return db.AutoMigrate(&Pet{})
}
