package main

import (
	"citra-wastra-be/config"
	"citra-wastra-be/handler"
	"citra-wastra-be/middleware"
	"citra-wastra-be/repository"
	"citra-wastra-be/routes"
	"citra-wastra-be/service"
	"citra-wastra-be/utils"
	"citra-wastra-be/worker"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {

	db := config.InitDB()
	rdb := config.InitRedis()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatalf("gemini initialization failed: %v", err)
	}
	defer client.Close()

	modelID := os.Getenv("GEMINI_MODEL_ID")
	narrator := utils.NewNarrator(client, modelID)
	uploader, err := utils.NewUploader()
	if err != nil {
		log.Fatalf("failed to initialize Cloudinary: %v", err)
	}
	classifier := utils.NewClassifier()

	// DI - repo
	userRepo := repository.NewUserRepository(db)
	batikRepo := repository.NewBatikRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	learningRepo := repository.NewLearningRepository(db)
	systemRepo := repository.NewSystemRepository(db)

	// redis
	cacheRepo := repository.NewBatikCacheRepository(rdb)
	queueRepo := repository.NewQueueRepository(rdb)
	gamificationRepo := repository.NewGamificationRepository(rdb)

	// DI - service
	gamificationService := service.NewGamificationService(badgeRepo, gamificationRepo, userRepo, rdb)
	authService := service.NewAuthService(userRepo)
	learningService := service.NewLearningService(learningRepo, queueRepo)
	adminService := service.NewAdminService(userRepo, learningRepo, batikRepo, badgeRepo, systemRepo, queueRepo)
	superAdminService := service.NewSuperAdminService(userRepo, systemRepo, queueRepo, learningRepo, badgeRepo)
	chatbotService := service.NewChatbotService(client, modelID)
	batikService := service.NewBatikService(
		batikRepo,
		cacheRepo,
		queueRepo,
		gamificationService,
		narrator,
		uploader,
		classifier)

	// worker
	xpWorker := worker.NewXPWorker(gamificationService, queueRepo)
	go xpWorker.Start(context.Background())

	// DI - handler
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authService)
	batikHandler := handler.NewBatikHandler(batikService)
	gamificationHandler := handler.NewGamificationHandler(gamificationService, gamificationRepo, userRepo)
	learningHandler := handler.NewLearningHandler(learningService)
	adminHandler := handler.NewAdminHandler(adminService)
	superAdminHandler := handler.NewSuperAdminHandler(superAdminService)
	chatbotHandler := handler.NewChatbotHandler(chatbotService)
	// router
	utils.StartKeepAlive()
	router := gin.Default()
	router.Use(middleware.StatsMiddleware(rdb))

	routes.SetupRoutes(router, authHandler, healthHandler, batikHandler, gamificationHandler, learningHandler, adminHandler, superAdminHandler, chatbotHandler, systemRepo)

	router.Run()
}
