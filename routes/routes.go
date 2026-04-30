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
	learningHandler *handler.LearningHandler,
	adminHandler *handler.AdminHandler,
) {

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		learning := api.Group("/learning")
		{
			learning.GET("/islands", learningHandler.GetIslands)
			learning.GET("/islands/:id/modules", learningHandler.GetModules)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{

			protected.GET("/leaderboard", gamificationHandler.GetLeaderboard)

			protected.GET("/user/me", authHandler.GetProfile)
			protected.POST("/batik/detect", batikHandler.Detect)
			protected.GET("/profile/badges", gamificationHandler.GetMyBadges)
			protected.POST("/profile/badges/equip", gamificationHandler.EquipBadge)
			protected.POST("/learning/levels/:id/complete", learningHandler.CompleteLevel)
			protected.POST("/learning/levels/:id/quiz", learningHandler.SubmitQuiz)

		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// CMS Levels
			admin.POST("/learning/levels", adminHandler.CreateLevel)
			admin.PUT("/learning/levels/:id", adminHandler.UpdateLevel)
			admin.DELETE("/learning/levels/:id", adminHandler.DeleteLevel)

			// CMS Questions
			admin.POST("/learning/questions", adminHandler.CreateQuestion)
			admin.PUT("/learning/questions/:id", adminHandler.UpdateQuestion)
			admin.DELETE("/learning/questions/:id", adminHandler.DeleteQuestion)

			// CMS Batik Catalog
			admin.GET("/batik/catalog", adminHandler.GetAllCatalog)
			admin.POST("/batik/catalog", adminHandler.CreateCatalog)
			admin.PUT("/batik/catalog/:id", adminHandler.UpdateCatalog)
			admin.DELETE("/batik/catalog/:id", adminHandler.DeleteCatalog)

			// CMS Badges
			admin.POST("/badges", adminHandler.CreateBadge)
			admin.PUT("/badges/:id", adminHandler.UpdateBadge)
			admin.DELETE("/badges/:id", adminHandler.DeleteBadge)

			// Monitoring
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.GET("/batik/logs", adminHandler.GetDetectionLogs)
			admin.GET("/health", adminHandler.GetSystemHealth)
		}
	}
}
