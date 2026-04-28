package routes

import (
	"citra-wastra-be/handler"
	"citra-wastra-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	batikHandler *handler.BatikHandler,
	gamificationHandler *handler.GamificationHandler,
) {

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/user/me", authHandler.GetProfile)
			protected.POST("/batik/detect", batikHandler.Detect)
			protected.GET("/profile/badges", gamificationHandler.GetMyBadges)
			protected.POST("/profile/badges/equip", gamificationHandler.EquipBadge)
		}
	}
}
