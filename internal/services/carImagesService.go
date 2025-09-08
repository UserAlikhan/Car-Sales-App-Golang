package services

import (
	"car_sales/internal/cache"
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"car_sales/internal/utils"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadCarPostImages(ctx *gin.Context, s3Conf *configs.S3Config, bucketName string, files []*multipart.FileHeader, carPostID uint) ([]string, error) {
	if len(files) == 0 {
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

		// create a path in s3 bucket where file is gonna be stored
		key := fmt.Sprintf("car_posts_photos/%d/%s", int(carPostID), fileHeader.Filename)

		// upload to s3
		_, err = utils.UploadToS3(
			ctx, s3Conf, &key, file,
			fileHeader.Header.Get("Content-Type"),
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to upload the image. ", err.Error())
		}

		// close the file after we are done
		file.Close()

		// if image was successfully uploaded to s3 create a database record
		carImage, err := repositories.CreateCarImage(key, uint(carPostID))
		if err != nil {
			return nil, err
		}

		// get signed url
		signedUrl, err := utils.GetSignedUrl(ctx, s3Conf, os.Getenv("BUCKET_NAME"), key, 24*time.Hour)
		if err != nil {
			return nil, fmt.Errorf("Unable to receive the url for the image. ", err.Error())
		}

		uploadedURLs = append(uploadedURLs, signedUrl)

		// store new image into the Redis cache
		cacheKey := fmt.Sprintf("carpost:%d:image:%d", carPostID, carImage.ID)
		cache.SetCache(cacheKey, signedUrl, 24*time.Hour)

		// After uploading all the images, refresh the carpost cache
		carPost, err := repositories.GetCarPostByIdWithPostImages(carPostID)
		if err == nil {
			marshaledCarPost, _ := json.Marshal(carPost)
			cache.SetCache(fmt.Sprintf("carpost:%d", carPostID), string(marshaledCarPost), 24*time.Hour)
		}
	}

	return uploadedURLs, nil
}

func GetCarPostImagesURLs(ctx *gin.Context, s3Conf *configs.S3Config, bucketName string, carPostID uint) ([]string, error) {
	// check if car post exists
	carPost, err := GetCarPostByIDWithoutImageURLs(carPostID)
	if err != nil {
		return nil, err
	}

	// slice for storing signedURLs
	var signedURLs []string

	// iterate through each post image
	for _, image := range carPost.PostImages {
		// get signed url
		signedUrl, err := utils.GetSignedUrl(ctx, s3Conf, s3Conf.BucketName, image.Path, 24*time.Hour)
		if err != nil {
			return nil, err
		}

		// save image to the slice
		signedURLs = append(signedURLs, signedUrl)
	}

	return signedURLs, nil
}

func GetCarImageByIDWithoutURL(ID int) (*models.CarImagesModel, error) {
	return repositories.GetCarImageByID(ID)
}

func DeleteCarImage(ctx *gin.Context, s3Conf *configs.S3Config, ID int) error {
	// check if image exists
	carImage, err := GetCarImageByIDWithoutURL(ID)
	if err != nil {
		return err
	}

	// delete image from the Redis cache
	cacheKey := fmt.Sprintf("carpost:%d:image:%d", carImage.CarPostID, ID)
	cache.DeleteCache(cacheKey)

	// delete image from s3 bucket
	err = utils.DeleteFromS3(ctx, s3Conf, carImage.Path)
	if err != nil {
		return err
	}

	// delete image record from database
	err = repositories.DeleteCarImageDBRecord(uint(ID))
	if err != nil {
		return err
	}

	// After deleting the images, refresh the carpost cache
	carPost, err := repositories.GetCarPostByIdWithPostImages(carImage.CarPostID)
	if err == nil {
		marshaledCarPost, _ := json.Marshal(carPost)
		cache.SetCache(fmt.Sprintf("carpost:%d", carImage.CarPostID), string(marshaledCarPost), 24*time.Hour)
	}

	return nil
}
