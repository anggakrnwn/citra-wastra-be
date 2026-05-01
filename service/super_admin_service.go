package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/models"
	"citra-wastra-be/repository"
	"citra-wastra-be/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type SuperAdminService interface {
	// Admin Management
	CreateAdmin(adminID, ip string, req dto.SuperAdminCreateAdminRequest) error
	UpdateUserStatus(adminID, userID, ip string, req dto.SuperAdminUpdateUserStatusRequest) error
	UpdateUserRole(adminID, userID, ip string, req dto.SuperAdminUpdateRoleRequest) error
	DeleteUser(adminID, userID, ip string) error

	// System Config
	UpdateConfig(adminID, ip string, req dto.SuperAdminConfigDTO) error
	GetAllConfigs() ([]models.SystemConfig, error)
	DeleteConfig(adminID, ip, key string) error

	// Audit Trail
	GetAuditLogs(page, limit int) ([]dto.SuperAdminAuditLogResponse, int64, error)
	LogActivity(adminID, action, target, details, ip string) error

	// Technical
	ClearXPQueue(ctx context.Context) error
	ResetTestData() error
}

type superAdminService struct {
	userRepo     repository.UserRepository
	systemRepo   repository.SystemRepository
	queueRepo    repository.QueueRepository
	learningRepo repository.LearningRepository
	badgeRepo    repository.BadgeRepository
}

func NewSuperAdminService(
	userRepo repository.UserRepository,
	systemRepo repository.SystemRepository,
	queueRepo repository.QueueRepository,
	learningRepo repository.LearningRepository,
	badgeRepo repository.BadgeRepository,
) SuperAdminService {
	return &superAdminService{userRepo, systemRepo, queueRepo, learningRepo, badgeRepo}
}

func (s *superAdminService) CreateAdmin(adminID, ip string, req dto.SuperAdminCreateAdminRequest) error {
	hashedPassword, _ := utils.HashPassword(req.Password)
	user := &models.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "admin",
	}
	err := s.userRepo.Create(user)
	if err == nil {
		s.LogActivity(adminID, "CREATE_ADMIN", req.Username, fmt.Sprintf("Created admin account for %s", req.Email), ip)
	}
	return err
}

func (s *superAdminService) UpdateUserStatus(adminID, userID, ip string, req dto.SuperAdminUpdateUserStatusRequest) error {
	err := s.userRepo.UpdateUserStatus(userID, req.IsBanned, req.IsActive)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_USER_STATUS", userID, fmt.Sprintf("Banned: %v, Active: %v", req.IsBanned, req.IsActive), ip)
	}
	return err
}

func (s *superAdminService) UpdateUserRole(adminID, userID, ip string, req dto.SuperAdminUpdateRoleRequest) error {
	err := s.userRepo.UpdateUserRole(userID, req.Role)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_USER_ROLE", userID, fmt.Sprintf("Role changed to %s", req.Role), ip)
	}
	return err
}

func (s *superAdminService) DeleteUser(adminID, userID, ip string) error {
	err := s.userRepo.DeleteUser(userID)
	if err == nil {
		s.LogActivity(adminID, "DELETE_USER", userID, "Permanently deleted user", ip)
	}
	return err
}

func (s *superAdminService) UpdateConfig(adminID, ip string, req dto.SuperAdminConfigDTO) error {
	config := &models.SystemConfig{
		Key:      req.Key,
		Value:    req.Value,
		Category: req.Category,
	}
	err := s.systemRepo.SetConfig(config)
	if err == nil {
		s.LogActivity(adminID, "UPDATE_CONFIG", req.Key, fmt.Sprintf("Value set to: %s", req.Value), ip)
	}
	return err
}

func (s *superAdminService) GetAllConfigs() ([]models.SystemConfig, error) {
	return s.systemRepo.GetAllConfigs()
}

func (s *superAdminService) DeleteConfig(adminID, ip, key string) error {
	err := s.systemRepo.DeleteConfig(key)
	if err == nil {
		s.LogActivity(adminID, "DELETE_CONFIG", key, "Config key deleted", ip)
	}
	return err
}

func (s *superAdminService) GetAuditLogs(page, limit int) ([]dto.SuperAdminAuditLogResponse, int64, error) {
	offset := (page - 1) * limit
	logs, total, err := s.systemRepo.GetAuditLogs(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.SuperAdminAuditLogResponse
	for _, l := range logs {
		res = append(res, dto.SuperAdminAuditLogResponse{
			ID:        l.ID,
			AdminName: l.Admin.Username,
			Action:    l.Action,
			Target:    l.Target,
			Details:   l.Details,
			IPAddress: l.IPAddress,
			CreatedAt: l.CreatedAt,
		})
	}
	return res, total, nil
}

func (s *superAdminService) LogActivity(adminID, action, target, details, ip string) error {
	log := &models.AuditLog{
		AdminID:   adminID,
		Action:    action,
		Target:    target,
		Details:   details,
		IPAddress: ip,
	}
	return s.systemRepo.CreateAuditLog(log)
}

func (s *superAdminService) ClearXPQueue(ctx context.Context) error {
	return s.queueRepo.ClearQueue(ctx)
}

func (s *superAdminService) ResetTestData() error {
	return s.systemRepo.ResetTestData()
}
