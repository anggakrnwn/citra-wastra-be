package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	sqlDB, _ := h.db.DB()
	dbStatus := "OK"
	if err := sqlDB.Ping(); err != nil {
		dbStatus = "DOWN"
	}

	c.JSON(http.StatusOK, gin.H{
		"app":     "citra-wastra-backend",
		"status":  "OK",
		"message": "API is running!",
		"details": gin.H{
			"database": dbStatus,
			"version":  "1.0.0",
			"time":     time.Now().Format(time.RFC3339),
		},
	})
}
