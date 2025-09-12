package model

import "gorm.io/gorm"

type Topup struct {
	BaseModel
	OrderId           int64 `gorm:"unique"`
	Type              int
	Email             string
	ShopifyCustomerId string `gorm:"size:64;not null;index" json:"shopifyCustomerId"`
}

// 创建表
func MigrateTopup(db *gorm.DB) error {
	return db.AutoMigrate(&Topup{})
}
