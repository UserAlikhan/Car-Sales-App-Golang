package handlers

import (
	"car_sales/internal/models"
	"car_sales/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCarBrandsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func GetCarBrandByIdHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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

		createdCarBrand, err := services.CreateCarBrand(carBrandData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, createdCarBrand)
	}
}

func UpdateCarBrandHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteCarBrandHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
