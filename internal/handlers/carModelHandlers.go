package handlers

import (
	"car_sales/internal/models"
	"car_sales/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCarModelHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var carModel *models.CarModelsModel

		if err := ctx.ShouldBindJSON(&carModel); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := services.CreateCarModel(carModel)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carModel)
	}
}

func GetAllCarModelsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		carModels, err := services.GetAllCarModels()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carModels)
	}
}

func GetCarModelsByCarBrandIDHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		carBrandID, err := strconv.Atoi(ctx.Param("carBrandID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
			return
		}

		carModels, err := services.GetCarModelsByCarBrandID(carBrandID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carModels)
	}
}
