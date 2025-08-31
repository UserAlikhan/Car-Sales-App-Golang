package handlers

import (
	"car_sales/internal/models"
	"car_sales/internal/services"
	"net/http"
	"strconv"

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
