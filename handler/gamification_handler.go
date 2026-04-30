package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/middleware"
	"citra-wastra-be/repository"
	"citra-wastra-be/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GamificationHandler struct {
	service  service.GamificationService
	gamiRepo repository.GamificationRepository
	userRepo repository.UserRepository
}

func NewGamificationHandler(s service.GamificationService, cr repository.GamificationRepository, ur repository.UserRepository) *GamificationHandler {
	return &GamificationHandler{s, cr, ur}
}

func (h *GamificationHandler) GetMyBadges(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)

	xp, _ := h.gamiRepo.GetUserXP(c.Request.Context(), userID)

	response, err := h.service.GetBadgeGallery(userID, int(xp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *GamificationHandler) EquipBadge(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)

	var req dto.EquipBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "id badge must be submitted",
		})
		return
	}

	if err := h.service.EquipBadge(userID, req.BadgeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successfully changed badge",
	})
}

func (h *GamificationHandler) GetLeaderboard(c *gin.Context) {
	scores, err := h.gamiRepo.GetLeaderboard(c.Request.Context(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var userIDs []string
	for _, s := range scores {
		userIDs = append(userIDs, s.Member.(string))
	}

	users, err := h.userRepo.GetUsersByIDs(userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to fetch user names",
		})
		return
	}

	userMap := make(map[string]string)
	for _, u := range users {
		userMap[u.ID] = u.Username
	}

	type LeaderboardItem struct {
		Rank     int     `json:"rank"`
		Username string  `json:"username"`
		XP       float64 `json:"xp"`
	}

	var leaderboard []LeaderboardItem
	for i, s := range scores {
		userID := s.Member.(string)
		name := userMap[userID]
		if name == "" {
			name = "Wastra Player"
		}

		leaderboard = append(leaderboard, LeaderboardItem{
			Rank:     i + 1,
			Username: name,
			XP:       s.Score,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    leaderboard,
	})
}
