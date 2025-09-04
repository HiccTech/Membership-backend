package router

import (
	"hiccpet/service/handler"
	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fmt"
	"hiccpet/service/config"
	"log"
)

func migrateDB(db *gorm.DB) {
	model.Migrate(db)
	model.MigrateStore(db)

	if err := model.MigrateCustomer(db); err != nil {
		log.Fatal("migrate customer error:", err)
	}

	if err := model.MigratePet(db); err != nil {
		log.Fatal("migrate pet error:", err)
	}

}

func runStatic(r *gin.Engine) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	r.Static("/static", filepath.Join(cwd, "uploadPetImgs"))
}

func SetupRouter() *gin.Engine {

	// 初始化配置
	config.LoadConfig()
	fmt.Println("当前环境 =", config.Cfg.ShopEnv)
	fmt.Println("店铺名称 =", config.Cfg.StoreName)
	fmt.Println("token =", config.Cfg.AccessToken)
	fmt.Println("数据库连接 =", config.Cfg.DbDSN)

	// 连接 MySQL
	dsn := config.Cfg.DbDSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	migrateDB(db)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	runStatic(r)

	// 公共接口
	r.POST("/register", func(c *gin.Context) { handler.Register(c, db) })
	r.POST("/login", func(c *gin.Context) { handler.Login(c, db) })

	// 定义 storefront 分组
	storefront := r.Group("/storefront")
	storefront.Use(middleware.ShopifySessionAuth())
	{
		storefront.GET("/getPetsByShopifyCustomerID", func(c *gin.Context) {
			handler.GetPetsByShopifyCustomerID(c, db)
		})

		storefront.POST("/addPet", func(c *gin.Context) {
			handler.AddPet(c, db)
		})

		storefront.POST("/updatePetById", func(c *gin.Context) {
			handler.UpdatePetById(c, db)
		})

		storefront.POST("/deletePetById", func(c *gin.Context) {
			handler.DeletePetById(c, db)
		})

		storefront.POST("/uploadPetAvatar", func(c *gin.Context) {
			handler.UploadPetAvatar(c, db)
		})
	}

	// 受保护接口
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	{

		auth.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is protected data"})
		})
		auth.POST("/addStore", func(c *gin.Context) { handler.AddStore(c, db) })
		auth.GET("/getStores", func(c *gin.Context) { handler.GetStores(c, db) })
	}

	return r
}
