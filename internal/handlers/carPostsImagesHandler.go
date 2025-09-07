package handlers

import (
	"car_sales/internal/configs"
	"car_sales/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// After user uploads new images, it reloads the page with car post
// this method uploads new images to s3 bucket and returns
// signed urls for newly created and old images
func UploadImagesHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get carPostID parameter passed from URL
		carPostID, err := strconv.Atoi(ctx.Param("carPostID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter."})
			return
		}

		// get all car post's images before uploding new ones
		carPostImagesURLs, err := services.GetCarPostImagesURLs(ctx, s3Conf, s3Conf.BucketName, uint(carPostID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// get files from the body
		files := ctx.Request.MultipartForm.File["photos"]

		// upload images and get urls back
		imageURLs, err := services.UploadCarPostImages(ctx, s3Conf, s3Conf.BucketName, files, uint(carPostID))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// append old images array with newly created urls
		carPostImagesURLs = append(carPostImagesURLs, imageURLs...)

		ctx.JSON(http.StatusCreated, gin.H{
			"imageURLs": carPostImagesURLs,
		})
	}
}

func DeleteCarPostImageHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ID parameter from endpoint's url
		ID, err := strconv.Atoi(ctx.Param("ID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
			return
		}

		// call DeleteCarImage service that will delete image from s3 and DB
		err = services.DeleteCarImage(ctx, s3Conf, ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Image was deleted successfully"})
	}
}
