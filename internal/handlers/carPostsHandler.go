package handlers

import (
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCarPostHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse form
		// if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form: " + err.Error()})
		// 	return
		// }

		var carPost models.CarPostsModel

		jsonStr := ctx.PostForm("car_post")
		if err := json.Unmarshal([]byte(jsonStr), &carPost); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car_post JSON."})
			return
		}

		if err := ctx.ShouldBindJSON(&carPost); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdCarPost, err := services.CreateCarPost(&carPost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// If car post was created proceed to photos
		// Get uploaded files for photos
		// files := ctx.Request.MultipartForm.File["photos"]

		// imageUrls, err := services.UploadCarPostImages(ctx, s3Conf, s3Conf.BucketName, files, createdCarPost.ID)
		// if err != nil {
		// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }

		ctx.JSON(http.StatusCreated, gin.H{
			"message":  "Car post was created succsesfully",
			"car_post": createdCarPost,
			// "image_urls": ,
		})
	}
}
