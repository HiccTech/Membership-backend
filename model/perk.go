package model

import "gorm.io/gorm"

type Perk struct {
	BaseModel
	ShopifyCustomerId string `gorm:"size:64;not null;index" json:"shopifyCustomerId"`
}

// 创建表
func MigratePerk(db *gorm.DB) error {
	return db.AutoMigrate(&Perk{})
}
