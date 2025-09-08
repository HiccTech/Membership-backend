package main

import (
	"hiccpet/service/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8081")
}
