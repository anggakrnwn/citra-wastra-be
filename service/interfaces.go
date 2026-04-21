package service

import "mime/multipart"

type NarratorInterface interface {
	GenerateNarration(label string) (string, error)
}

type ImageUploader interface {
	UploadImage(file *multipart.FileHeader, folder string) (string, error)
}

type BatikClassifier interface {
	ClassifyBatik(imageURL string) (string, float64, string, error)
}
