package services

import (
	"mime/multipart"
	config "project/Config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadToCloudinary(ctx *gin.Context, file multipart.File, filePath string) (string, error) {
	// ctx := context.Background()
	cld, err := config.ConfigCloudinary()
	if err != nil {
		return "", err
	}
	uploadParams := uploader.UploadParams{
		PublicID: filePath,
	}
	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	imageUrl := result.SecureURL
	return imageUrl, nil
}
