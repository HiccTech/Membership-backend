package config

import (
	"os"
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
	"testShop": {
		"StoreName":     "Test Store",
		"CountryCode":   "test",
		"StoreDomain":   "test-store-hicc1.myshopify.com",
		"AccessToken":   "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",
		"Admin":         "https://admin.shopify.com/store/test-store-hicc1",
		"DbDSN":         "root:Neo123456!@tcp(127.0.0.1:3306)/membership_test?charset=utf8mb4&parseTime=True&loc=Local",
		"WebhookSecret": "3d7242eb2f79c5055faf0addd8642252799ccb1a2d750cbfd59b4e272245e623",
	},
	"sgShop": {
		"StoreName":     "Prod Store",
		"CountryCode":   "prod",
		"StoreDomain":   "prod-store.myshopify.com",
		"AccessToken":   "shpat_prod_xxx",
		"Admin":         "https://admin.shopify.com/store/prod-store",
		"DbDSN":         "root:password@tcp(127.0.0.1:3306)/membership_prod?charset=utf8mb4&parseTime=True&loc=Local",
		"WebhookSecret": "",
	},
}

func LoadConfig() {
	env := os.Getenv("SHOP_ENV")
	if env == "" {
		env = "testShop"
	}

	Cfg = &Config{
		ShopEnv:       env,
		StoreName:     shopCfg[env]["StoreName"],
		StoreDomain:   shopCfg[env]["StoreDomain"],
		AccessToken:   shopCfg[env]["AccessToken"],
		DbDSN:         shopCfg[env]["DbDSN"],
		WebhookSecret: shopCfg[env]["WebhookSecret"],
	}

}
