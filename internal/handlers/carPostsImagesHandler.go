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
	return func(c *gin.Context) {
		key := c.Query("path") // e.g. "car_posts_photos/9/Ford_photo.png"

		err := utils.DeleteFromS3(c, s3Conf, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image was deleted successfully"})
	}
}
