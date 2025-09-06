package services

import (
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"car_sales/internal/utils"
	"os"
	"strings"
	"time"

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
	carPost, err := repositories.GetCarPostByIdWithPostImages(ID)
	if err != nil {
		return err
	}
	// If there are images saved proceed
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

func GetCarPostByID(ctx *gin.Context, s3Conf *configs.S3Config, ID uint) (*models.CarPostsModel, []string, error) {
	// get a car post with preloaded images
	carPost, err := repositories.GetCarPostByIdWithPostImages(ID)
	if err != nil {
		return nil, nil, err
	}

	// array for storing signed urls
	signedURLs := make([]string, 0, len(carPost.PostImages))

	// get a prefix
	parts := strings.Split(carPost.PostImages[0].Path, "/")
	parts = parts[:len(parts)-1]
	prefix := strings.Join(parts, "/")

	// get list of all present images in s3 bucket for the specific prefix
	listOutput, err := s3Conf.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &s3Conf.BucketName,
		Prefix: &prefix,
	})
	if err != nil {
		return nil, nil, err
	}

	// variable for storing keys from listoutput
	var listOutputKeys []*string
	// get all the keys from listoutput
	for _, output := range listOutput.Contents {
		listOutputKeys = append(listOutputKeys, output.Key)
	}

	// iterate through all images for the post
	// and get signed url for post images
	for _, image := range carPost.PostImages {
		// get signed urls
		signedURL, err := utils.GetSignedUrl(ctx, s3Conf, os.Getenv("BUCKET_NAME"), image.Path, 24*time.Hour)
		if err != nil {
			// call FindStringKeyStringValue method that looks for an image in s3
			// if indeed there is no image in s3 delete the record from DB
			if !utils.FindStringKeyStringValue(listOutputKeys, &image.Path) {
				// delete the post image record
				repositories.DeleteCarImageDBRecord(image.ID)
			} else {
				continue
			}
		}
		// if everything is good append an array
		signedURLs = append(signedURLs, signedURL)
	}

	return carPost, signedURLs, nil
}

func GetNumberOfCarPosts() int {
	return repositories.CountCarPostsTotalRows()
}

func GetCarPostsWithPagination(context *gin.Context, s3Conf *configs.S3Config, limit int, page int) ([]*models.CarPostsModel, error) {
	// calculate offset
	offset := (page - 1) * limit

	// get paginated data
	carPosts, err := repositories.GetCarPostsWithPagination(limit, offset)
	if err != nil {
		return nil, err
	}

	// this is made for the website, where there are car posts cards with one photo
	// we do not need to make signed urls for all of the photos in carPosts, we need only one
	// that will be shown on the card. Later if user opens a post we will make another call
	// that will load all the photos related to the car post
	for _, carPost := range carPosts {
		for _, image := range carPost.PostImages {
			// get presigned url
			signedURL, err := utils.GetSignedUrl(context, s3Conf, s3Conf.BucketName, image.Path, 24*time.Hour)
			if err != nil {
				continue
			}
			// store newly generated url
			carPost.CardPhotoURL = signedURL
		}
	}

	return carPosts, err
}
