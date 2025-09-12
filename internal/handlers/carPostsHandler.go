package handlers

import (
	"car_sales/internal/configs"
	"car_sales/internal/models"
	"car_sales/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCarPostHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var carPost models.CarPostsModel

		// get json data from car_post parameter
		jsonStr := ctx.PostForm("car_post")
		if err := json.Unmarshal([]byte(jsonStr), &carPost); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car_post JSON."})
			return
		}

		// get userID parameter from Authorization middleware
		userID := ctx.GetInt("userID")

		// car post's seller ID should be the same as
		// userID from the jwt token
		carPost.SellerID = uint(userID)

		// create car post in the database
		createdCarPost, err := services.CreateCarPost(&carPost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// If car post was created proceed to photos
		// Get uploaded photos
		files := ctx.Request.MultipartForm.File["photos"]

		// call service layer that will proceed images upload
		imageUrls, err := services.UploadCarPostImages(ctx, s3Conf, s3Conf.BucketName, files, createdCarPost.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"car_post":   createdCarPost,
			"image_urls": imageUrls,
		})
	}
}

func GetAllUsersCarPostsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := strconv.Atoi(ctx.Param("userId"))
		if err != nil || userId == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalida credentials."})
			return
		}

		carPosts, err := services.GetAllUsersCarPosts(uint(userId))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carPosts)
	}
}

func DeleteCarPostHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ID parameter
		ID, err := strconv.Atoi(ctx.Param("ID"))
		if err != nil || ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials."})
			return
		}

		// call service layer that will proceed car post deletion
		err = services.DeleteCarPost(ctx, s3Conf, uint(ID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Car Post was deleted successfully"})
	}
}

func GetCarPostByIDHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ID parameter
		ID, err := strconv.Atoi(ctx.Param("ID"))
		if err != nil || ID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials."})
			return
		}

		carPost, image_urls, err := services.GetCarPostByID(ctx, s3Conf, uint(ID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(
			http.StatusOK,
			gin.H{
				"carPost":    carPost,
				"image_urls": image_urls,
			},
		)
	}
}

func GetCarPostsWithPaginationHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, err := strconv.Atoi(ctx.Query("Limit"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No 'Limit' parameter was found."})
			return
		}

		page, err := strconv.Atoi(ctx.Query("Page"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No 'Page' parameter was found."})
			return
		}

		carPosts, err := services.GetCarPostsWithPagination(ctx, s3Conf, limit, page)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carPosts)
	}
}

func UpdateCarPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID, err := strconv.Atoi(ctx.Param("ID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
			return
		}

		var carPost *models.CarPostsModel

		if err := ctx.ShouldBindJSON(&carPost); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		carPost.ID = uint(ID)

		// call service to update car post record
		err = services.UpdateCarPost(carPost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carPost)
	}
}
