package router

import (
	"hiccpet/service/handler"
	"hiccpet/service/middleware"
	"hiccpet/service/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fmt"

	"hiccpet/service/config"
)

func SetupRouter() *gin.Engine {
	// 初始化配置
	config.LoadConfig()
	fmt.Println("当前环境 =", config.Cfg.ShopEnv)
	fmt.Println("店铺名称 =", config.Cfg.StoreName)
	fmt.Println("token =", config.Cfg.AccessToken)
	fmt.Println("数据库连接 =", config.Cfg.DBDsn)

	// 连接 MySQL
	dsn := config.Cfg.DBDsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	model.Migrate(db)
	model.MigrateStore(db)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	// 公共接口
	r.POST("/register", func(c *gin.Context) { handler.Register(c, db) })
	r.POST("/login", func(c *gin.Context) { handler.Login(c, db) })

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
