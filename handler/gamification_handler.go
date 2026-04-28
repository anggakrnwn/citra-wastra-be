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
}

func NewGamificationHandler(s service.GamificationService, cr repository.GamificationRepository) *GamificationHandler {
	return &GamificationHandler{s, cr}
}

func (h *GamificationHandler) GetMyBadges(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)

	xp, _ := h.gamiRepo.GetUserXP(c.Request.Context(), userID)

	response, err := h.service.GetBadgeGallery(userID, int(xp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
			"error": "id badge must be submitted",
		})
		return
	}

	if err := h.service.EquipBadge(userID, req.BadgeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successfully changed badge"})
}
