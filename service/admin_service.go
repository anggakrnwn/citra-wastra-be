package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AdminService interface {
	// CMS - Levels
	CreateLevel(adminID, ip string, req dto.AdminLevelRequest) error
	UpdateLevel(adminID, id, ip string, req dto.AdminLevelRequest) error
	DeleteLevel(adminID, id, ip string) error

	// CMS - Modules
	CreateModule(adminID, ip string, req dto.AdminModuleRequest) error

	// CMS - Questions
	CreateQuestion(adminID, ip string, req dto.AdminQuestionRequest) error
	UpdateQuestion(adminID, id, ip string, req dto.AdminQuestionRequest) error
	DeleteQuestion(adminID, id, ip string) error

	// CMS - Batik Catalog
	CreateCatalog(adminID, ip string, req dto.AdminCatalogRequest) error
	UpdateCatalog(adminID, id, ip string, req dto.AdminCatalogRequest) error
	DeleteCatalog(adminID, id, ip string) error
	GetAllCatalog() ([]models.BatikCatalog, error)

	// CMS - Badges
	CreateBadge(adminID, ip string, req dto.AdminBadgeRequest) error
	UpdateBadge(adminID, id, ip string, req dto.AdminBadgeRequest) error
	DeleteBadge(adminID, id, ip string) error

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
	systemRepo   repository.SystemRepository
	queueRepo    repository.QueueRepository
}

func NewAdminService(
	userRepo repository.UserRepository,
	learningRepo repository.LearningRepository,
	batikRepo repository.BatikRepository,
	badgeRepo repository.BadgeRepository,
	systemRepo repository.SystemRepository,
	queueRepo repository.QueueRepository,
) AdminService {
	return &adminService{userRepo, learningRepo, batikRepo, badgeRepo, systemRepo, queueRepo}
}

func (s *adminService) LogActivity(adminID, action, target, details, ip string) {
	log := &models.AuditLog{
		AdminID:   adminID,
		Action:    action,
		Target:    target,
		Details:   details,
		IPAddress: ip,
	}
	s.systemRepo.CreateAuditLog(log)
}

func (s *adminService) CreateLevel(adminID, ip string, req dto.AdminLevelRequest) error {
	level := &models.Level{
		ID:       uuid.NewString(),
		ModuleID: req.ModuleID,
		Title:    req.Title,
		Content:  req.Content,
		Order:    req.Order,
		XPReward: req.XPReward,
	}
	err := s.learningRepo.CreateLevel(level)
	if err == nil {
		s.LogActivity(adminID, "CREATE_LEVEL", level.ID, fmt.Sprintf("Title: %s", req.Title), ip)
	}
	return err
}

func (s *adminService) UpdateLevel(adminID, id, ip string, req dto.AdminLevelRequest) error {
	level := &models.Level{
		ID:       id,
		ModuleID: req.ModuleID,
		Title:    req.Title,
		Content:  req.Content,
		Order:    req.Order,
		XPReward: req.XPReward,
	}
	err := s.learningRepo.UpdateLevel(level)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_LEVEL", id, fmt.Sprintf("Title: %s", req.Title), ip)
	}
	return err
}

func (s *adminService) DeleteLevel(adminID, id, ip string) error {
	err := s.learningRepo.DeleteLevel(id)
	if err == nil {
		s.LogActivity(adminID, "DELETE_LEVEL", id, "Deleted level", ip)
	}
	return err
}

func (s *adminService) CreateModule(adminID, ip string, req dto.AdminModuleRequest) error {
	module := &models.Module{
		ID:          uuid.NewString(),
		IslandID:    req.IslandID,
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.learningRepo.CreateModule(module)
	if err == nil {
		s.LogActivity(adminID, "CREATE_MODULE", module.ID, fmt.Sprintf("Name: %s", req.Name), ip)
	}
	return err
}

func (s *adminService) CreateQuestion(adminID, ip string, req dto.AdminQuestionRequest) error {
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
	err := s.learningRepo.CreateQuestion(q)
	if err == nil {
		s.LogActivity(adminID, "CREATE_QUESTION", q.ID, fmt.Sprintf("LevelID: %s", req.LevelID), ip)
	}
	return err
}

func (s *adminService) UpdateQuestion(adminID, id, ip string, req dto.AdminQuestionRequest) error {
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
	err := s.learningRepo.UpdateQuestion(q)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_QUESTION", id, "Updated question content", ip)
	}
	return err
}

func (s *adminService) DeleteQuestion(adminID, id, ip string) error {
	err := s.learningRepo.DeleteQuestion(id)
	if err == nil {
		s.LogActivity(adminID, "DELETE_QUESTION", id, "Deleted question", ip)
	}
	return err
}

func (s *adminService) CreateCatalog(adminID, ip string, req dto.AdminCatalogRequest) error {
	catalog := &models.BatikCatalog{
		ID:        uuid.NewString(),
		Name:      req.Name,
		History:   req.History,
		Origin:    req.Origin,
		ImageURL:  req.ImageURL,
		CreatedBy: &adminID,
	}
	err := s.batikRepo.CreateCatalog(catalog)
	if err == nil {
		s.LogActivity(adminID, "CREATE_CATALOG", catalog.ID, fmt.Sprintf("Name: %s", req.Name), ip)
	}
	return err
}

func (s *adminService) UpdateCatalog(adminID, id, ip string, req dto.AdminCatalogRequest) error {
	catalog := &models.BatikCatalog{
		ID:       id,
		Name:     req.Name,
		History:  req.History,
		Origin:   req.Origin,
		ImageURL: req.ImageURL,
	}
	err := s.batikRepo.UpdateCatalog(catalog)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_CATALOG", id, fmt.Sprintf("Name: %s", req.Name), ip)
	}
	return err
}

func (s *adminService) DeleteCatalog(adminID, id, ip string) error {
	err := s.batikRepo.DeleteCatalog(id)
	if err == nil {
		s.LogActivity(adminID, "DELETE_CATALOG", id, "Deleted catalog item", ip)
	}
	return err
}

func (s *adminService) GetAllCatalog() ([]models.BatikCatalog, error) {
	return s.batikRepo.GetAllCatalog(), nil
}

func (s *adminService) CreateBadge(adminID, ip string, req dto.AdminBadgeRequest) error {
	badge := &models.Badge{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		MinXP:       req.MinXP,
		ImageURL:    req.ImageURL,
	}
	err := s.badgeRepo.CreateBadge(badge)
	if err == nil {
		s.LogActivity(adminID, "CREATE_BADGE", badge.ID, fmt.Sprintf("Name: %s", req.Name), ip)
	}
	return err
}

func (s *adminService) UpdateBadge(adminID, id, ip string, req dto.AdminBadgeRequest) error {
	badge := &models.Badge{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		MinXP:       req.MinXP,
		ImageURL:    req.ImageURL,
	}
	err := s.badgeRepo.UpdateBadge(badge)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_BADGE", id, fmt.Sprintf("Name: %s", req.Name), ip)
	}
	return err
}

func (s *adminService) DeleteBadge(adminID, id, ip string) error {
	err := s.badgeRepo.DeleteBadge(id)
	if err == nil {
		s.LogActivity(adminID, "DELETE_BADGE", id, "Deleted badge", ip)
	}
	return err
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
			ID:         u.ID,
			Username:   u.Username,
			Email:      u.Email,
			XP:         u.XP,
			Role:       u.Role,
			IsVerified: u.IsVerified,
			IsBanned:   u.IsBanned,
			IsActive:   u.IsActive,
			CreatedAt:  u.CreatedAt,
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
