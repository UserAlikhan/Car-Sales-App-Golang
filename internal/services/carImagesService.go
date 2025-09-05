package services

import (
	"car_sales/internal/configs"
	"car_sales/internal/repositories"
	"car_sales/internal/utils"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadCarPostImages(ctx *gin.Context, s3Conf *configs.S3Config, bucketName string, files []*multipart.FileHeader, carPostID uint) ([]string, error) {
	if len(files) < 0 {
		return nil, fmt.Errorf("No images uploaded.")
	}

	if len(files) > 10 {
		return nil, fmt.Errorf("User allowed to upload up to 10 photos at the same time.")
	}

	uploadedURLs := make([]string, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("Failed to open file: ", err.Error())
		}

		defer file.Close()

		key := fmt.Sprintf("car_posts_photos/%d/%s", int(carPostID), fileHeader.Filename)

		// upload to s3
		_, err = utils.UploadToS3(
			ctx, s3Conf, &key, file,
			fileHeader.Header.Get("Content-Type"),
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to upload the image. ", err.Error())
		}

		// if image was successfully uploaded to s3 create a database record
		err = repositories.CreateCarImage(key, uint(carPostID))

		// get signed url
		signedUrl, err := utils.GetSignedUrl(ctx, s3Conf, os.Getenv("BUCKET_NAME"), key, 24*time.Hour)
		if err != nil {
			return nil, fmt.Errorf("Unable to receive the url for the image. ", err.Error())
		}

		uploadedURLs = append(uploadedURLs, signedUrl)
	}

	return uploadedURLs, nil
}
