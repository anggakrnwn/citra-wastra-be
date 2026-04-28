package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	MinConfidenceThreshold = 0.65
	CacheDuration          = 24 * time.Hour
	RedisPrefix            = "batik_cache"
)

type BatikService interface {
	DetectBatik(ctx context.Context, userID string, req dto.UploadBatikRequest) (dto.BatikResponse, error)
}

type batikService struct {
	repo          repository.BatikRepository
	cacheRepo     repository.BatikCacheRepository
	queueRepo     repository.QueueRepository
	narrator      NarratorInterface
	uploader      ImageUploader
	classifier    BatikClassifier
	gammification GamificationService
}

func NewBatikService(
	repo repository.BatikRepository,
	cacheRepo repository.BatikCacheRepository,
	queueRepo repository.QueueRepository,
	gammification GamificationService,
	narrator NarratorInterface,
	uploader ImageUploader,
	classifier BatikClassifier,
) BatikService {
	return &batikService{
		repo:          repo,
		cacheRepo:     cacheRepo,
		queueRepo:     queueRepo,
		gammification: gammification,
		narrator:      narrator,
		uploader:      uploader,
		classifier:    classifier,
	}
}

func (s *batikService) DetectBatik(ctx context.Context, userID string, req dto.UploadBatikRequest) (dto.BatikResponse, error) {

	fileHash, err := utils.CalculateFileHash(req.File)
	if err != nil {
		return dto.BatikResponse{}, fmt.Errorf("failed to calculate hash: %v", err)
	}

	cacheKey := RedisPrefix + ":" + fileHash
	if cachedData, err := s.cacheRepo.Get(ctx, cacheKey); err == nil {
		var cachedResponse dto.BatikResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			log.Printf("cache hit for key: %s", cacheKey)
			return cachedResponse, nil
		}
	}

	folderName := os.Getenv("CLOUDINARY_FOLDER_BATIK")
	if folderName == "" {
		folderName = "default-batik"
	}
	imageURL, err := s.uploader.UploadImage(ctx, req.File, folderName)
	if err != nil {
		return dto.BatikResponse{}, fmt.Errorf("cloudinary upload failed: %v", err)
	}

	label, confidence, philosophy, err := s.classifier.ClassifyBatik(ctx, imageURL)
	if err != nil {
		return dto.BatikResponse{}, fmt.Errorf("ai classification failed: %v", err)
	}

	var narrativeInput string
	var finalLabel string

	if confidence >= MinConfidenceThreshold {
		narrativeInput = label
		finalLabel = label
	} else {
		narrativeInput = fmt.Sprintf("Object with visual characteristics resembling %s (Confidence: %.2f)", label, confidence)
		finalLabel = "unidentified motif"
	}

	finalPhilosophy, err := s.narrator.GenerateNarration(narrativeInput)
	if err != nil || finalPhilosophy == "" {
		log.Printf("narrator error: %v", err)

		finalPhilosophy = philosophy
		fmt.Println("using static philosophy due to narrator error")
	}

	optimizedURL := utils.OptimizeImageURL(imageURL)

	batik := models.Batik{
		UserID:             userID,
		ImageURL:           optimizedURL,
		Label:              finalLabel,
		Confidence:         confidence,
		PhilosophyNarrator: finalPhilosophy,
		Philosophy:         philosophy,
	}

	if err := s.repo.Create(&batik); err != nil {
		return dto.BatikResponse{}, fmt.Errorf("failed to save batik data: %v", err)
	}

	response := dto.BatikResponse{
		ID:         batik.ID,
		ImageURL:   batik.ImageURL,
		Label:      finalLabel,
		Confidence: batik.Confidence,
		Philosophy: finalPhilosophy,
		CreatedAt:  batik.CreatedAt,
	}

	if err := s.cacheRepo.Set(ctx, cacheKey, response, CacheDuration); err != nil {
		log.Printf("warning: failed to set cache: %v", err)
	}

	payload := dto.XPJobPayload{
		JobID:  uuid.NewString(),
		UserID: userID,
		XPGain: 100,
		Retry:  0,
	}

	if err := s.queueRepo.EnqueueXPJob(ctx, payload); err != nil {
		log.Printf("failed to increment XP: %v", err)
	}

	return response, nil
}
