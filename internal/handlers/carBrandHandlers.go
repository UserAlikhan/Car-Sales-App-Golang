package handlers

import (
	"car_sales/internal/configs"
	"car_sales/internal/database"
	"car_sales/internal/models"
	"car_sales/internal/services"
	"car_sales/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func GetAllCarBrandsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		carBrands, err := services.GetAllCarBrands()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carBrands)
	}
}

func GetCarBrandByIdHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// call the service
		carBrand, err := services.GetCarBrandById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carBrand)
	}
}

func CreateCarBrandHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var carBrandData *models.CarBrandsModel
		// Get car brand data and store it to the variable
		if err := ctx.ShouldBindJSON(&carBrandData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := services.CreateCarBrand(carBrandData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// return car brand data if successfull
		ctx.JSON(http.StatusCreated, carBrandData)
	}
}

func UpdateCarBrandHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var carBrand *models.CarBrandsModel
		if err := ctx.ShouldBindJSON(&carBrand); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		carBrand.ID = uint(ID)

		err = services.UpdateCarBrand(carBrand)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carBrand)
	}
}

func DeleteCarBrandHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ID, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = services.DeleteCarBrand(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, "Car Brand was deleted successfully")
	}
}

func CreateCarBrandWithModelsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var carBrand models.CarBrandsModel
		if err := ctx.ShouldBindJSON(&carBrand); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := services.CreateCarBrandWithModels(&carBrand)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carBrand)
	}
}

func UploadLogoHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID param
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse Uploaded file
		file, header, err := ctx.Request.FormFile("logo")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		// check if car brand with the given id exists
		carBrand, err := services.GetCarBrandById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Car Brand with given credentials was not found."})
			return
		}

		// filename
		key := fmt.Sprintf("car_brands/%d/%s", id, header.Filename)

		// upload to S3
		_, err = s3Conf.Uploader.Upload(ctx, &s3.PutObjectInput{
			Bucket:      &s3Conf.BucketName,
			Key:         &key,
			Body:        file,
			ContentType: aws.String(header.Header.Get("Content-Type")),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// save logo image to carBrand variable
		carBrand.LogoImage = key
		database.DB.Save(&carBrand)

		signedURL, err := utils.GetSignedUrl(ctx, s3Conf, s3Conf.BucketName, key, 15*time.Minute)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":   "Logo was uploaded successfully",
			"logo_url":  signedURL,
			"car_brand": carBrand,
		})
	}
}
