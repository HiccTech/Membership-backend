package model

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	StoreName   string
	CountryCode string `gorm:"unique"`
	StoreDomain string
	AccessToken string
	Admin       string
}

// 创建表
func MigrateStore(db *gorm.DB) error {
	return db.AutoMigrate(&Store{})
}
