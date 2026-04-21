package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"fmt"
	"log"
	"os"
)

const MinConfidenceThreshold = 0.65

type BatikService interface {
	DetectBatik(userID string, req dto.UploadBatikRequest) (dto.BatikResponse, error)
}

type batikService struct {
	repo       repository.BatikRepository
	narrator   NarratorInterface
	uploader   ImageUploader
	classifier BatikClassifier
}

func NewBatikService(
	repo repository.BatikRepository,
	narrator NarratorInterface,
	uploader ImageUploader,
	classifier BatikClassifier,
) BatikService {
	return &batikService{
		repo:       repo,
		narrator:   narrator,
		uploader:   uploader,
		classifier: classifier,
	}
}

func (s *batikService) DetectBatik(userID string, req dto.UploadBatikRequest) (dto.BatikResponse, error) {

	folderName := os.Getenv("CLOUDINARY_FOLDER_BATIK")
	if folderName == "" {
		folderName = "default-batik"
	}
	imageURL, err := s.uploader.UploadImage(req.File, folderName)
	if err != nil {
		return dto.BatikResponse{}, fmt.Errorf("cloudinary upload failed: %v", err)
	}

	label, confidence, philosophy, err := s.classifier.ClassifyBatik(imageURL)
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

	return dto.BatikResponse{
		ID:         batik.ID,
		ImageURL:   batik.ImageURL,
		Label:      finalLabel,
		Confidence: batik.Confidence,
		Philosophy: finalPhilosophy,
	}, nil
}
