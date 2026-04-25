package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/middleware"
	"citra-wastra-be/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BatikHandler struct {
	service service.BatikService
}

func NewBatikHandler(s service.BatikService) *BatikHandler {
	return &BatikHandler{s}
}

func (h *BatikHandler) Detect(c *gin.Context) {
	val, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := val.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal identity error",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image file is required",
		})
		return
	}

	req := dto.UploadBatikRequest{
		File: file,
	}

	result, err := h.service.DetectBatik(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrImageUpload):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrClassification):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			log.Printf("detect batik error (user_id=%s): %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to process batik detection",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "batik detected successfully",
		"data":    result,
	})
}
