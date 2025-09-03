package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

// 创建表
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
