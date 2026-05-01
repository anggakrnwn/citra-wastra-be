package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuperAdminHandler struct {
	service service.SuperAdminService
}

func NewSuperAdminHandler(s service.SuperAdminService) *SuperAdminHandler {
	return &SuperAdminHandler{s}
}

// Admin Management
func (h *SuperAdminHandler) CreateAdmin(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.SuperAdminCreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.CreateAdmin(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "admin account created successfully"})
}

func (h *SuperAdminHandler) UpdateUserStatus(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	var req dto.SuperAdminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateUserStatus(adminID, id, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user status updated successfully"})
}

func (h *SuperAdminHandler) UpdateUserRole(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	var req dto.SuperAdminUpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateUserRole(adminID, id, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user role updated successfully"})
}

func (h *SuperAdminHandler) DeleteUser(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	if err := h.service.DeleteUser(adminID, id, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "user permanently deleted successfully"})
}

// System Config
func (h *SuperAdminHandler) UpdateConfig(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.SuperAdminConfigDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateConfig(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "system config updated successfully"})
}

func (h *SuperAdminHandler) GetAllConfigs(c *gin.Context) {
	res, err := h.service.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

// Audit Logs
func (h *SuperAdminHandler) GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, total, err := h.service.GetAuditLogs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// Technical
func (h *SuperAdminHandler) ClearXPQueue(c *gin.Context) {
	if err := h.service.ClearXPQueue(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "XP queue cleared successfully"})
}

func (h *SuperAdminHandler) ResetTestData(c *gin.Context) {
	if err := h.service.ResetTestData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Test data reset successfully"})
}
