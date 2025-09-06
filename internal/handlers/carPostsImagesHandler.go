package handlers

import (
	"car_sales/internal/configs"
	"car_sales/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadImagesHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteCarPostImageHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Query("path") // e.g. "car_posts_photos/9/Ford_photo.png"

		err := utils.DeleteFromS3(ctx, s3Conf, key)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Image was deleted successfully"})
	}
}
