package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	service service.ChatbotService
}

func NewChatbotHandler(s service.ChatbotService) *ChatbotHandler {
	return &ChatbotHandler{s}
}

func (h *ChatbotHandler) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "message is required"})
		return
	}

	res, err := h.service.Chat(c.Request.Context(), req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}
