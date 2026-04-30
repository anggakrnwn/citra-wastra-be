package dto

import "time"

type SuperAdminCreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type SuperAdminUpdateUserStatusRequest struct {
	IsBanned bool `json:"is_banned"`
	IsActive bool `json:"is_active"`
}

type SuperAdminUpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type SuperAdminConfigDTO struct {
	Key      string `json:"key" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Category string `json:"category"`
}

type SuperAdminAuditLogResponse struct {
	ID        string    `json:"id"`
	AdminName string    `json:"admin_name"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
