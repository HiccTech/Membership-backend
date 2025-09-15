package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	ShopEnv       string
	StoreName     string
	StoreDomain   string
	AccessToken   string
	DbDSN         string
	WebhookSecret string
}

var shopCfg = map[string]map[string]string{
	"dev": {
		"StoreName":   "Test Store",
		"CountryCode": "test",
		"StoreDomain": "test-store-hicc1.myshopify.com",
		"AccessToken": "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",
		"Admin":       "https://admin.shopify.com/store/test-store-hicc1",
		// "DbDSN":         "root:Neo123456!@tcp(127.0.0.1:3306)/membership_test?charset=utf8mb4&parseTime=True&loc=Local",
		"DbUser":        os.Getenv("DB_USER"),
		"DbPassword":    os.Getenv("DB_PASSWORD"),
		"DbHost":        os.Getenv("DB_HOST"),
		"DbPort":        os.Getenv("DB_PORT"),
		"DbName":        os.Getenv("DB_NAME"),
		"WebhookSecret": "3d7242eb2f79c5055faf0addd8642252799ccb1a2d750cbfd59b4e272245e623",
	},
	"testShop": {
		"StoreName":   "Test Store",
		"CountryCode": "test",
		"StoreDomain": "test-store-hicc1.myshopify.com",
		"AccessToken": "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",
		"Admin":       "https://admin.shopify.com/store/test-store-hicc1",
		// "DbDSN":         "root:Neo123456!@tcp(127.0.0.1:3306)/membership_test?charset=utf8mb4&parseTime=True&loc=Local",
		"DbUser":        os.Getenv("DB_USER"),
		"DbPassword":    os.Getenv("DB_PASSWORD"),
		"DbHost":        os.Getenv("DB_HOST"),
		"DbPort":        os.Getenv("DB_PORT"),
		"DbName":        os.Getenv("DB_NAME"),
		"WebhookSecret": "3d7242eb2f79c5055faf0addd8642252799ccb1a2d750cbfd59b4e272245e623",
	},
	"sgShop": {
		"StoreName":   "Prod Store",
		"CountryCode": "prod",
		"StoreDomain": "prod-store.myshopify.com",
		"AccessToken": "shpat_prod_xxx",
		"Admin":       "https://admin.shopify.com/store/prod-store",
		// "DbDSN":         "root:password@tcp(127.0.0.1:3306)/membership_prod?charset=utf8mb4&parseTime=True&loc=Local",
		"DbUser":        os.Getenv("DB_USER"),
		"DbPassword":    os.Getenv("DB_PASSWORD"),
		"DbHost":        os.Getenv("DB_HOST"),
		"DbPort":        os.Getenv("DB_PORT"),
		"DbName":        os.Getenv("DB_NAME"),
		"WebhookSecret": "",
	},
}

func LoadConfig() {
	env := os.Getenv("SHOP_ENV")
	fmt.Println("Shop Env:", env)
	if env == "" {
		env = "testShop"
	}

	godotenv.Load(".env." + env)

	// 再读取 DB 配置
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	fmt.Println("DB DSN:", dbDSN)

	Cfg = &Config{
		ShopEnv:       env,
		StoreName:     shopCfg[env]["StoreName"],
		StoreDomain:   shopCfg[env]["StoreDomain"],
		AccessToken:   shopCfg[env]["AccessToken"],
		DbDSN:         dbDSN,
		WebhookSecret: shopCfg[env]["WebhookSecret"],
	}

}
