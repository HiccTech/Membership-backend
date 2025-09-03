package config

import (
	"os"

	"github.com/joho/godotenv"
)

// 全局变量
var Cfg *Config

type Config struct {
	ShopEnv     string
	StoreName   string
	AccessToken string
	DBDsn       string
}

func LoadConfig() {
	env := os.Getenv("SHOP_ENV")
	if env == "" {
		env = "testShop"
	}

	fileName := ".env." + env

	err := godotenv.Load(fileName)

	if err != nil {
		panic("加载配置文件失败")
	}

	Cfg = &Config{
		ShopEnv:     os.Getenv("SHOP_ENV"),
		StoreName:   os.Getenv("StoreName"),
		AccessToken: os.Getenv("AccessToken"),
		DBDsn:       os.Getenv("DB_DSN"),
	}

}
