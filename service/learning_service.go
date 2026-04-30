package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"context"
	"fmt"
)

type LearningService interface {
	GetAllIslands() ([]models.Island, error)
	GetModulesByIsland(islandID string) ([]models.Module, error)
	CompleteLevel(userID string, levelID string) error
	GetLevelDetail(levelID string) (models.Level, []models.Question, error)
	SubmitQuiz(userID string, levelID string, isCheating bool, score int) error
}

type learningService struct {
	learningRepo repository.LearningRepository
	queueRepo    repository.QueueRepository
}

func NewLearningService(lr repository.LearningRepository, qr repository.QueueRepository) LearningService {
	return &learningService{lr, qr}
}

func (s *learningService) GetAllIslands() ([]models.Island, error) {
	return s.learningRepo.GetIslands()
}

func (s *learningService) GetModulesByIsland(islandID string) ([]models.Module, error) {
	return s.learningRepo.GetModulesByIsland(islandID)
}

func (s *learningService) CompleteLevel(userID string, levelID string) error {
	level, err := s.learningRepo.GetLevelByID(levelID)
	if err != nil {
		return fmt.Errorf("level not found: %v", err)
	}

	progress := &models.UserLevelProgress{
		UserID:  userID,
		LevelID: levelID,
	}
	if err := s.learningRepo.SaveProgress(progress); err != nil {
		return fmt.Errorf("failed to save progress: %v", err)
	}

	payload := dto.XPJobPayload{
		JobID:  fmt.Sprintf("progress-%s-%s", userID, levelID),
		UserID: userID,
		XPGain: level.XPReward,
	}

	return s.queueRepo.EnqueueXPJob(context.Background(), payload)
}

func (s *learningService) GetLevelDetail(levelID string) (models.Level, []models.Question, error) {
	level, err := s.learningRepo.GetLevelByID(levelID)
	if err != nil {
		return models.Level{}, nil, err
	}
	questions, err := s.learningRepo.GetQuestionsByLevel(levelID)
	return level, questions, err
}

func (s *learningService) SubmitQuiz(userID string, levelID string, isCheating bool, score int) error {
	_, err := s.learningRepo.GetProgress(userID, levelID)
	if err != nil {
		return fmt.Errorf("you must complete the level content before taking the quiz")
	}

	level, err := s.learningRepo.GetLevelByID(levelID)
	if err != nil {
		return err
	}

	questions, _ := s.learningRepo.GetQuestionsByLevel(levelID)
	totalQuestions := len(questions)

	xpGain := 0
	if totalQuestions > 0 {
		xpGain = (score * level.XPReward) / totalQuestions
	} else {
		xpGain = level.XPReward
	}

	if isCheating {
		xpGain -= 25
	}

	payload := dto.XPJobPayload{
		JobID:  fmt.Sprintf("quiz-%s-%s", userID, levelID),
		UserID: userID,
		XPGain: xpGain,
	}

	return s.queueRepo.EnqueueXPJob(context.Background(), payload)
}
