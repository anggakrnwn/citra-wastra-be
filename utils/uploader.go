package utils

import (
	"context"
	"mime/multipart"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type cloudinaryUploader struct {
	client *cloudinary.Cloudinary
}

func NewUploader() (*cloudinaryUploader, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return nil, err
	}
	return &cloudinaryUploader{client: cld}, nil

}

func (u *cloudinaryUploader) UploadImage(fileHeader *multipart.FileHeader, folder string) (string, error) {

	ctx := context.Background()
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploadResult, err := u.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil

}

func OptimizeImageURL(url string) string {
	return strings.Replace(url, "/upload/", "/upload/q_auto,f_auto/", 1)
}
