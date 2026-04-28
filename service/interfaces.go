package service

import (
	"context"
	"mime/multipart"
)

type NarratorInterface interface {
	GenerateNarration(label string) (string, error)
}

type ImageUploader interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
}

type BatikClassifier interface {
	ClassifyBatik(ctx context.Context, imageURL string) (string, float64, string, error)
}
