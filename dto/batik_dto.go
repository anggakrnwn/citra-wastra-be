package dto

import (
	"mime/multipart"
	"time"
)

type UploadBatikRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type BatikResponse struct {
	ID         string    `json:"id"`
	ImageURL   string    `json:"image_url"`
	Label      string    `json:"label"`
	Confidence float64   `json:"confidence"`
	Philosophy string    `json:"philosophy,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
