package router

import (
	"fmt"
	"hiccpet/service/handler"
	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"hiccpet/service/service"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	migrateDB(db)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	runStatic(r)

	sseApp := service.NewSSEServer()
	r.GET("/sse", sseApp.Handler)

	// 模拟每秒推送一次消息给 customer_id=123
	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		msg := fmt.Sprintf("tick: %s", time.Now().Format("15:04:05"))
	// 		sseApp.PushToClient("123", msg)
	// 		sseApp.PushToClient("bob123", "hixx")
	// 		sseApp.PushToClient("neo123", "hixx--neo")
	// 	}
	// }()

	// 公共接口
	r.POST("/register", func(c *gin.Context) { handler.Register(c, db) })
	r.POST("/login", func(c *gin.Context) { handler.Login(c, db) })

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
	}

	r.POST("/webhook/orders", middleware.ShopifyWebhookAuth(), func(c *gin.Context) {
		b, _ := io.ReadAll(c.Request.Body)
		fmt.Println(string(b), " -----------------------")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": string(b)})
	})

	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	{

		auth.POST("/addStore", func(c *gin.Context) { handler.AddStore(c, db) })
		auth.GET("/getStores", func(c *gin.Context) { handler.GetStores(c, db) })
	}

	return r
}
