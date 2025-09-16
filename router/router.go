package router

import (
	"hiccpet/service/email"
	"hiccpet/service/handler"
	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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

	if err := model.MigrateTopup(db); err != nil {
		log.Fatal("migrate topup error:", err)
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

	// 连接 MySQL
	dsn := config.Cfg.DbDSN
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ { // 最多尝试 10 次
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("Waiting for database to be ready...", err)
		time.Sleep(2 * time.Second) // 每 2 秒尝试一次
	}
	if err != nil {
		log.Fatal("failed to connect to database after retries:", err)
	}

	migrateDB(db)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	runStatic(r)

	// sseApp := service.NewSSEServer()
	// r.GET("/sse/club", sseApp.Handler)

	// 公共接口
	r.POST("/register", func(c *gin.Context) { handler.Register(c, db) })
	r.POST("/login", func(c *gin.Context) { handler.Login(c, db) })

	r.GET("/sendEmail", func(c *gin.Context) {
		email.SendClubEmail(email.EmailData{To: "812284688@qq.com", Subject: "Test Email", DiscountCodes: []email.DiscountCode{
			{Title: "Test Code 1", Code: "TEST1", Period: "2024/01/01 - 2024/12/31"},
			{Title: "Test Code 2", Code: "TEST2", Period: "2024/01/01 - 2024/12/31"},
			{Title: "Test Code 2", Code: "TEST2", Period: "2024/01/01 - 2024/12/31"},
		}})
		c.JSON(200, gin.H{"message": "uccessful"})
	})

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

		storefront.GET("/getStoreCreditBalance", func(c *gin.Context) {
			handler.GetStoreCreditBalance(c, db)
		})

		storefront.POST("/getCodeDiscountNodes", func(c *gin.Context) {
			handler.GetCodeDiscountNodes(c, db)
		})

		storefront.GET("topupCount", func(c *gin.Context) {
			handler.TopupCount(c, db)
		})

	}

	r.POST("/webhook/orders", middleware.ShopifyWebhookAuth(), func(c *gin.Context) {
		handler.HandleTopUp(c, db)
	})

	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	{

		// auth.POST("/addStore", func(c *gin.Context) { handler.AddStore(c, db) })
		// auth.GET("/getStores", func(c *gin.Context) { handler.GetStores(c, db) })
	}

	return r
}
