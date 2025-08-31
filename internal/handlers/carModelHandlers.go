package handlers

import (
	"car_sales/internal/models"
	"car_sales/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCarModel() gin.HandlerFunc {
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
