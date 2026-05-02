package routes

import (
	"citra-wastra-be/handler"
	"citra-wastra-be/middleware"
	"citra-wastra-be/repository"

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
	superAdminHandler *handler.SuperAdminHandler,
	chatbotHandler *handler.ChatbotHandler,
	systemRepo repository.SystemRepository,
) {

	r.GET("/health", healthHandler.Check)

	api := r.Group("/api/v1")
	api.Use(middleware.MaintenanceMiddleware(systemRepo))
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)
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
			protected.GET("/learning/levels/:id", learningHandler.GetLevelDetail)
			protected.POST("/learning/levels/:id/complete", learningHandler.CompleteLevel)
			protected.POST("/learning/levels/:id/quiz", learningHandler.SubmitQuiz)

			// Chatbot
			protected.POST("/chat", chatbotHandler.Chat)
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

		superAdmin := api.Group("/super-admin")
		superAdmin.Use(middleware.AuthMiddleware(), middleware.SuperAdminMiddleware())
		{
			// Admin & User Management
			superAdmin.POST("/admins", superAdminHandler.CreateAdmin)
			superAdmin.PUT("/users/:id/status", superAdminHandler.UpdateUserStatus)
			superAdmin.PUT("/users/:id/role", superAdminHandler.UpdateUserRole)
			superAdmin.PUT("/users/:id/reset-password", superAdminHandler.ResetPassword)
			superAdmin.DELETE("/users/:id", superAdminHandler.DeleteUser)

			// System Configuration
			superAdmin.GET("/config", superAdminHandler.GetAllConfigs)
			superAdmin.PUT("/config", superAdminHandler.UpdateConfig)

			// Audit Logs
			superAdmin.GET("/audit-logs", superAdminHandler.GetAuditLogs)

			// Technical
			superAdmin.POST("/technical/clear-queue", superAdminHandler.ClearXPQueue)
			superAdmin.POST("/technical/reset-data", superAdminHandler.ResetTestData)
		}
	}
}
