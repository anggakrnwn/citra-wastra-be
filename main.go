package main

import (
	"citra-wastra-be/config"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitDB()

	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "API running",
		})
	})

	router.Run()
}
