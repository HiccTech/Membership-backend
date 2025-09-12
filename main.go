package main

import (
	"hiccpet/service/router"
	"os"
)

func main() {
	r := router.SetupRouter()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // 默认端口
	}
	r.Run(":" + port)

}
