package services

import (
	"car_sales/internal/cache"
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/repositories"
	"car_sales/internal/utils"
	"encoding/json"
	"fmt"
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

func GetCarPostByIDWithoutImageURLs(ID uint) (*models.CarPostsModel, error) {
	return repositories.GetCarPostById(ID)
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
	err = repositories.DeleteCarPost(ID)
	if err != nil {
		return err
	}

	// Eager Invalidation
	// After we deleted a car post we need to delete the cache
	// because it is not up-to-date
	// Here we don’t know exact limit/page user used, so best is to delete all carposts:* keys
	// TODO: After I implement car filtration, delete car posts without filters
	// and delete car posts with filters for exact the same car type

	// Starts from the begining and fetchs 100 keys per iteration
	iter := configs.RedisClient.Scan(ctx, 0, "carposts:*", 100).Iterator()
	for iter.Next(ctx) {
		_ = cache.DeleteCache(iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func GetCarPostByID(ctx *gin.Context, s3Conf *configs.S3Config, ID uint) (*models.CarPostsModel, []string, error) {
	// 1. Try to get car post from Redis cache first

	// car post key for accessing car post from the Redis cache
	carPostCacheKey := fmt.Sprintf("carpost:%d", ID)

	var carPost *models.CarPostsModel

	// try to get cache from Redis
	if cachedData, err := cache.GetCache(carPostCacheKey); err == nil && cachedData != "" {
		// unmarshal json int the variable
		if err := json.Unmarshal([]byte(cachedData), &carPost); err != nil {
			return nil, nil, err
		}

		// Get image urls from Redis cache
		imageURLs := []string{}

		// Images are saved individually, so search the images and put them into the slice
		iter := configs.RedisClient.Scan(ctx, 0, fmt.Sprintf("carpost:%d:image:*", ID), 0).Iterator()
		for iter.Next(ctx) {
			// get image value from the Redis cache
			url, err := cache.GetCache(iter.Val())
			if err != nil {
				return nil, nil, err
			}

			// append the slice
			imageURLs = append(imageURLs, url)
		}
		// if there is no error and len of images retrieved from the Redis cache equal
		// to the number of database return the data
		if err := iter.Err(); err == nil {
			return carPost, imageURLs, nil
		}
	}

	// 2. If there is no data in Redis cache fetch data from Database
	fmt.Println("GETTING DATA FROM DATABASE")
	// get a car post with preloaded images
	carPost, err := repositories.GetCarPostByIdWithPostImages(ID)
	if err != nil {
		return nil, nil, err
	}

	// array for storing signed urls
	signedURLs := make([]string, 0, len(carPost.PostImages))

	// if there are some post images, we need to get signed url for them
	if len(carPost.PostImages) > 0 {
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
				// convert listOutputKeys to be able to use it with Find function
				listOutputKeysConverted := utils.ConvertStringPointerArrayToStringArray(listOutputKeys)
				// call FindStringKeyStringValue method that looks for an image in s3
				// if indeed there is no image in s3 delete the record from DB
				if !utils.Find(listOutputKeysConverted, image.Path) {
					// delete the post image record
					repositories.DeleteCarImageDBRecord(image.ID)
				} else {
					continue
				}
			}

			// Save image URL to the Redis cache
			cache.SetCache(fmt.Sprintf("carpost:%d:image:%d", ID, image.ID), signedURL, 24*time.Hour)

			// if everything is good append an array
			signedURLs = append(signedURLs, signedURL)
		}
	}

	// Marshal CarPost
	marsheledCarPost, _ := json.Marshal(carPost)
	// Save car post to the Redis cache
	cache.SetCache(fmt.Sprintf("carpost:%d", ID), string(marsheledCarPost), 24*time.Hour)

	return carPost, signedURLs, nil
}

func GetNumberOfCarPosts() int {
	return repositories.CountCarPostsTotalRows()
}

func GetCarPostsWithPagination(context *gin.Context, s3Conf *configs.S3Config, limit int, page int) ([]*models.CarPostsModel, error) {
	// key to access redis cache data
	cacheKey := fmt.Sprintf("carposts:limit:%d:page:%d", limit, page)

	// 1. Try to get data from the Redis cache first
	if cached, err := cache.GetCache(cacheKey); err == nil {
		var carPosts []*models.CarPostsModel
		if err := json.Unmarshal([]byte(cached), &carPosts); err != nil {
			return nil, err
		}

		return carPosts, nil
	}

	// 2. If there is not data in Redis cache fetch from DB
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

	// Save data to the Redis cache
	response, _ := json.Marshal(carPosts)
	cache.SetCache(cacheKey, string(response), 24*time.Hour)

	return carPosts, err
}
