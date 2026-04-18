package routes

import (
	"citra-wastra-be/handler"
	"citra-wastra-be/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler, healthHandler *handler.HealthHandler) {

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/user")
		protected.Use(middleware.AuthMiddleware())
		{
			// test dummy
			protected.GET("/ping", func(ctx *gin.Context) {
				userID := ctx.MustGet(middleware.UserIDKey).(string)
				role := ctx.MustGet(middleware.RoleKey).(string)

				ctx.JSON(200, gin.H{
					"message": "Pong! Token kamu valid.",
					"data": gin.H{
						"user_id": userID,
						"role":    role,
					},
				})
			})
		}
	}
}
