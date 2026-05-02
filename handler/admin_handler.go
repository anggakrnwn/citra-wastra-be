package handler

import (
	"citra-wastra-be/dto"
	"citra-wastra-be/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service service.AdminService
}

func NewAdminHandler(s service.AdminService) *AdminHandler {
	return &AdminHandler{s}
}

// CMS - Levels
func (h *AdminHandler) CreateLevel(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.AdminLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.CreateLevel(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "level created successfully"})
}

func (h *AdminHandler) UpdateLevel(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	var req dto.AdminLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateLevel(adminID, id, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "level updated successfully"})
}

func (h *AdminHandler) DeleteLevel(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	if err := h.service.DeleteLevel(adminID, id, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "level deleted successfully"})
}

// CMS - Modules
func (h *AdminHandler) CreateModule(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.AdminModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.CreateModule(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "module created successfully"})
}

// CMS - Questions
func (h *AdminHandler) CreateQuestion(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.AdminQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.CreateQuestion(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "question created successfully"})
}

func (h *AdminHandler) UpdateQuestion(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	var req dto.AdminQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateQuestion(adminID, id, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "question updated successfully"})
}

func (h *AdminHandler) DeleteQuestion(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	if err := h.service.DeleteQuestion(adminID, id, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "question deleted successfully"})
}

// CMS - Batik Catalog
func (h *AdminHandler) CreateCatalog(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.AdminCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.CreateCatalog(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "catalog asset created successfully"})
}

func (h *AdminHandler) UpdateCatalog(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	var req dto.AdminCatalogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateCatalog(adminID, id, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "catalog asset updated successfully"})
}

func (h *AdminHandler) DeleteCatalog(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	if err := h.service.DeleteCatalog(adminID, id, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "catalog asset deleted successfully"})
}

func (h *AdminHandler) GetAllCatalog(c *gin.Context) {
	res, err := h.service.GetAllCatalog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

// CMS - Badges
func (h *AdminHandler) CreateBadge(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()

	var req dto.AdminBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.CreateBadge(adminID, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "badge created successfully"})
}

func (h *AdminHandler) UpdateBadge(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	var req dto.AdminBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.service.UpdateBadge(adminID, id, ip, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "badge updated successfully"})
}

func (h *AdminHandler) DeleteBadge(c *gin.Context) {
	adminID := c.GetString("user_id")
	ip := c.ClientIP()
	id := c.Param("id")

	if err := h.service.DeleteBadge(adminID, id, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "badge deleted successfully"})
}

// Monitoring
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, total, err := h.service.GetAllUsers(page, limit)
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

func (h *AdminHandler) GetDetectionLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, total, err := h.service.GetDetectionLogs(page, limit)
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

func (h *AdminHandler) GetSystemHealth(c *gin.Context) {
	health, err := h.service.GetSystemHealth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": health})
}

func (h *AdminHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "image is required"})
		return
	}

	url, err := h.service.UploadImage(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "url": url})
}
