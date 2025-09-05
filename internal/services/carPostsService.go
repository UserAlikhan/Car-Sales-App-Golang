package services

import (
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"car_sales/internal/utils"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func CreateCarPost(carPost *models.CarPostsModel) (*models.CarPostsModel, error) {
	return repositories.CreateCarPost(carPost)
}

func GetAllUsersCarPosts(userId uint) ([]*models.CarPostsModel, error) {
	return repositories.GetAllUsersCarPosts(userId)
}

func DeleteCarPost(ctx *gin.Context, s3Conf *configs.S3Config, ID uint) error {
	// Check if car post exists
	carPost, err := repositories.GetCarPostById(ID)
	if err != nil {
		return err
	}

	if len(carPost.PostImages) > 0 {
		// get prefix
		parts := strings.Split(carPost.PostImages[0].Path, "/")
		parts = parts[:len(parts)-1]
		// join back into the string
		prefix := strings.Join(parts, "/")

		// List of all objects under the prefix
		listOutput, err := s3Conf.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: &s3Conf.BucketName,
			Prefix: &prefix,
		})
		if err != nil {
			return err
		}

		for _, img := range listOutput.Contents {
			if img.Key != nil {
				// delete images from s3 bucket
				if err := utils.DeleteFromS3(ctx, s3Conf, *img.Key); err != nil {
					return err
				}
			}
		}
	}

	// call delete car post
	return repositories.DeleteCarPost(ID)
}
