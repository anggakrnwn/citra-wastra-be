package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"context"

	"github.com/google/uuid"
)

type AdminService interface {
	// CMS - Levels
	CreateLevel(req dto.AdminLevelRequest) error
	UpdateLevel(id string, req dto.AdminLevelRequest) error
	DeleteLevel(id string) error

	// CMS - Questions
	CreateQuestion(req dto.AdminQuestionRequest) error
	UpdateQuestion(id string, req dto.AdminQuestionRequest) error
	DeleteQuestion(id string) error

	// CMS - Batik Catalog
	CreateCatalog(adminID string, req dto.AdminCatalogRequest) error
	UpdateCatalog(id string, req dto.AdminCatalogRequest) error
	DeleteCatalog(id string) error
	GetAllCatalog() ([]models.BatikCatalog, error)

	// CMS - Badges
	CreateBadge(req dto.AdminBadgeRequest) error
	UpdateBadge(id string, req dto.AdminBadgeRequest) error
	DeleteBadge(id string) error

	// Monitoring
	GetAllUsers(page, limit int) ([]dto.AdminUserResponse, int64, error)
	GetDetectionLogs(page, limit int) ([]models.Batik, int64, error)
	GetSystemHealth(ctx context.Context) (dto.AdminSystemHealthResponse, error)
}

type adminService struct {
	userRepo     repository.UserRepository
	learningRepo repository.LearningRepository
	batikRepo    repository.BatikRepository
	badgeRepo    repository.BadgeRepository
	queueRepo    repository.QueueRepository
}

func NewAdminService(
	userRepo repository.UserRepository,
	learningRepo repository.LearningRepository,
	batikRepo repository.BatikRepository,
	badgeRepo repository.BadgeRepository,
	queueRepo repository.QueueRepository,
) AdminService {
	return &adminService{userRepo, learningRepo, batikRepo, badgeRepo, queueRepo}
}

func (s *adminService) CreateLevel(req dto.AdminLevelRequest) error {
	level := &models.Level{
		ID:       uuid.NewString(),
		ModuleID: req.ModuleID,
		Title:    req.Title,
		Content:  req.Content,
		Order:    req.Order,
		XPReward: req.XPReward,
	}
	return s.learningRepo.CreateLevel(level)
}

func (s *adminService) UpdateLevel(id string, req dto.AdminLevelRequest) error {
	level := &models.Level{
		ID:       id,
		ModuleID: req.ModuleID,
		Title:    req.Title,
		Content:  req.Content,
		Order:    req.Order,
		XPReward: req.XPReward,
	}
	return s.learningRepo.UpdateLevel(level)
}

func (s *adminService) DeleteLevel(id string) error {
	return s.learningRepo.DeleteLevel(id)
}

func (s *adminService) CreateQuestion(req dto.AdminQuestionRequest) error {
	q := &models.Question{
		ID:       uuid.NewString(),
		LevelID:  req.LevelID,
		Question: req.Question,
		OptionA:  req.OptionA,
		OptionB:  req.OptionB,
		OptionC:  req.OptionC,
		OptionD:  req.OptionD,
		Correct:  req.Correct,
	}
	return s.learningRepo.CreateQuestion(q)
}

func (s *adminService) UpdateQuestion(id string, req dto.AdminQuestionRequest) error {
	q := &models.Question{
		ID:       id,
		LevelID:  req.LevelID,
		Question: req.Question,
		OptionA:  req.OptionA,
		OptionB:  req.OptionB,
		OptionC:  req.OptionC,
		OptionD:  req.OptionD,
		Correct:  req.Correct,
	}
	return s.learningRepo.UpdateQuestion(q)
}

func (s *adminService) DeleteQuestion(id string) error {
	return s.learningRepo.DeleteQuestion(id)
}

func (s *adminService) CreateCatalog(adminID string, req dto.AdminCatalogRequest) error {
	catalog := &models.BatikCatalog{
		ID:        uuid.NewString(),
		Name:      req.Name,
		History:   req.History,
		Origin:    req.Origin,
		ImageURL:  req.ImageURL,
		CreatedBy: &adminID,
	}
	return s.batikRepo.CreateCatalog(catalog)
}

func (s *adminService) UpdateCatalog(id string, req dto.AdminCatalogRequest) error {
	catalog := &models.BatikCatalog{
		ID:       id,
		Name:     req.Name,
		History:  req.History,
		Origin:   req.Origin,
		ImageURL: req.ImageURL,
	}
	return s.batikRepo.UpdateCatalog(catalog)
}

func (s *adminService) DeleteCatalog(id string) error {
	return s.batikRepo.DeleteCatalog(id)
}

func (s *adminService) GetAllCatalog() ([]models.BatikCatalog, error) {
	return s.batikRepo.GetAllCatalog(), nil
}

func (s *adminService) CreateBadge(req dto.AdminBadgeRequest) error {
	badge := &models.Badge{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		MinXP:       req.MinXP,
		ImageURL:    req.ImageURL,
	}
	return s.badgeRepo.CreateBadge(badge)
}

func (s *adminService) UpdateBadge(id string, req dto.AdminBadgeRequest) error {
	badge := &models.Badge{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		MinXP:       req.MinXP,
		ImageURL:    req.ImageURL,
	}
	return s.badgeRepo.UpdateBadge(badge)
}

func (s *adminService) DeleteBadge(id string) error {
	return s.badgeRepo.DeleteBadge(id)
}

func (s *adminService) GetAllUsers(page, limit int) ([]dto.AdminUserResponse, int64, error) {
	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAllUsers(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.AdminUserResponse
	for _, u := range users {
		res = append(res, dto.AdminUserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			XP:        u.XP,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}
	return res, total, nil
}

func (s *adminService) GetDetectionLogs(page, limit int) ([]models.Batik, int64, error) {
	offset := (page - 1) * limit
	return s.batikRepo.GetAllLogs(limit, offset)
}

func (s *adminService) GetSystemHealth(ctx context.Context) (dto.AdminSystemHealthResponse, error) {
	qLen, _ := s.queueRepo.GetQueueLength(ctx)
	stats, _ := s.queueRepo.GetAPIStats(ctx)

	return dto.AdminSystemHealthResponse{
		WorkerStatus: "running",
		QueueLength:  qLen,
		Cloudinary:   "active",
		APIStats:     stats,
	}, nil
}
