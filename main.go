package main

import (
	"citra-wastra-be/config"
	"citra-wastra-be/handler"
	"citra-wastra-be/repository"
	"citra-wastra-be/routes"
	"citra-wastra-be/service"

	"github.com/gin-gonic/gin"
)

func main() {

	db := config.InitDB()

	// DI
	healthHandler := handler.NewHealthHandler(db)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	routes.SetupRoutes(router, authHandler, healthHandler)

	router.Run()
}
