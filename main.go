package main

import (
	"citra-wastra-be/config"
	"citra-wastra-be/handler"
	"citra-wastra-be/repository"
	"citra-wastra-be/routes"
	"citra-wastra-be/service"
	"citra-wastra-be/utils"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {

	db := config.InitDB()

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

	// DI
	healthHandler := handler.NewHealthHandler(db)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	batikRepo := repository.NewBatikRepository(db)
	batikService := service.NewBatikService(batikRepo, narrator, uploader, classifier)
	batikHandler := handler.NewBatikHandler(batikService)

	utils.StartKeepAlive()

	router := gin.Default()
	routes.SetupRoutes(router, authHandler, healthHandler, batikHandler)

	router.Run()
}
