package handler

import (
	"citra-wastra-be/middleware"
	"citra-wastra-be/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LearningHandler struct {
	learningService service.LearningService
}

func NewLearningHandler(s service.LearningService) *LearningHandler {
	return &LearningHandler{learningService: s}
}

func (h *LearningHandler) GetIslands(c *gin.Context) {
	islands, err := h.learningService.GetAllIslands()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    islands,
	})
}

func (h *LearningHandler) GetModules(c *gin.Context) {
	islandID := c.Param("id")
	userID, exists := c.Get(middleware.UserIDKey)

	var modules interface{}
	var err error

	if exists {
		modules, err = h.learningService.GetIslandModulesWithProgress(userID.(string), islandID)
	} else {
		modules, err = h.learningService.GetModulesByIsland(islandID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    modules,
	})
}

func (h *LearningHandler) CompleteLevel(c *gin.Context) {
	levelID := c.Param("id")

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	err := h.learningService.CompleteLevel(userID.(string), levelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "level completed success",
	})
}

func (h *LearningHandler) GetLevelDetail(c *gin.Context) {
	levelID := c.Param("id")
	level, questions, err := h.learningService.GetLevelDetail(levelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "level not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"level":     level,
		"questions": questions,
	})
}

func (h *LearningHandler) SubmitQuiz(c *gin.Context) {
	levelID := c.Param("id")
	var req struct {
		IsCheating bool `json:"is_cheating"`
		Score      int  `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request",
		})
		return
	}

	userID, _ := c.Get(middleware.UserIDKey)
	err := h.learningService.SubmitQuiz(userID.(string), levelID, req.IsCheating, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	message := "Quiz submitted!"
	if req.IsCheating {
		message = "help used. xp deducted"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}
