package service

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/repository"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type GamificationService interface {
	SyncUserBadges(userID string, currentXP int) error
	GetBadgeGallery(userID string, currentXP int) (dto.UserProfileBadgeResponse, error)
	EquipBadge(userID string, badgeID string) error
	ProcessXPAndBadges(payload dto.XPJobPayload) error
}

type gamificationService struct {
	badgeRepo repository.BadgeRepository
	gamiRepo  repository.GamificationRepository
	userRepo  repository.UserRepository
	rdb       *redis.Client
}

func NewGamificationService(badgeRepo repository.BadgeRepository, gamiRepo repository.GamificationRepository, userRepo repository.UserRepository, rdb *redis.Client) GamificationService {
	return &gamificationService{badgeRepo, gamiRepo, userRepo, rdb}
}

func (s *gamificationService) SyncUserBadges(userID string, currentXP int) error {
	allBadges, err := s.badgeRepo.GetAllBadges()
	if err != nil {
		return err
	}

	ownedIDs, err := s.badgeRepo.GetUnlockedBadgeIDs(userID)
	if err != nil {
		return err
	}

	ownedMap := make(map[string]bool)
	for _, id := range ownedIDs {
		ownedMap[id] = true
	}

	for _, b := range allBadges {
		if currentXP >= b.MinXP && !ownedMap[b.ID] {
			_ = s.badgeRepo.UnlockBadge(userID, b.ID)
		}
	}
	return nil
}

func (s *gamificationService) GetBadgeGallery(userID string, currentXP int) (dto.UserProfileBadgeResponse, error) {
	_ = s.SyncUserBadges(userID, currentXP)

	allBadges, err := s.badgeRepo.GetAllBadges()
	if err != nil {
		return dto.UserProfileBadgeResponse{}, err
	}

	ownedIDs, _ := s.badgeRepo.GetUnlockedBadgeIDs(userID)
	equipped, _ := s.badgeRepo.GetEquippedBadge(userID)

	ownedMap := make(map[string]bool)
	for _, id := range ownedIDs {
		ownedMap[id] = true
	}

	var badgeResponses []dto.BadgeResponse
	for _, b := range allBadges {

		isEquipped := false
		if equipped.BadgeID != "" && equipped.BadgeID == b.ID {
			isEquipped = true
		}

		badgeResponses = append(badgeResponses, dto.BadgeResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			ImageURL:    b.ImageURL,
			MinXP:       b.MinXP,
			IsUnlocked:  ownedMap[b.ID],
			IsEquipped:  isEquipped,
		})
	}

	return dto.UserProfileBadgeResponse{
		TotalXP: currentXP,
		Badges:  badgeResponses,
	}, nil
}

func (s *gamificationService) EquipBadge(userID string, badgeID string) error {
	ctx := context.Background()
	xp, _ := s.gamiRepo.GetUserXP(ctx, userID)
	_ = s.SyncUserBadges(userID, int(xp))

	ownedIDs, err := s.badgeRepo.GetUnlockedBadgeIDs(userID)
	if err != nil {
		return fmt.Errorf("failed to check ownership: %v", err)
	}

	found := false
	for _, id := range ownedIDs {
		if id == badgeID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("badge not unlocked yet")
	}

	return s.badgeRepo.EquipBadge(userID, badgeID)
}

func (s *gamificationService) ProcessXPAndBadges(payload dto.XPJobPayload) error {
	ctx := context.Background()

	key := "processed_job:" + payload.JobID

	err := s.rdb.SetArgs(ctx, key, true, redis.SetArgs{
		Mode: "NX",
		TTL:  24 * time.Hour,
	}).Err()

	if err == redis.Nil {
		log.Printf("job already processed: %s", payload.JobID)
		return nil
	} else if err != nil {
		return err
	}

	newXP, err := s.gamiRepo.IncrementUserXP(ctx, payload.UserID, float64(payload.XPGain))
	if err != nil {
		return err
	}

	if newXP < 0 {
		log.Printf("xp minus detected for user %s (%f). resetting to 0.", payload.UserID, newXP)

		if err := s.gamiRepo.SetUserXP(ctx, payload.UserID, 0); err != nil {
			return err
		}
		newXP = 0
	}

	if err := s.userRepo.UpdateTotalXP(payload.UserID, int(newXP)); err != nil {
		log.Printf("failed sync XP to DB for user %s: %v", payload.UserID, err)
	}

	return s.SyncUserBadges(payload.UserID, int(newXP))
}
