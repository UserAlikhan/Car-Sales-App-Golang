package handlers

import (
	"car_sales/internal/models"
	"car_sales/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userData *models.UsersModel

		// store response into the variable
		if err := ctx.ShouldBindJSON(&userData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// call service to create the new user
		user, err := services.CreateUser(userData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// return response if successfull
		ctx.JSON(http.StatusOK, user)
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginData *models.LoginDataModel

		// store the incoming data to the variable
		if err := ctx.ShouldBindJSON(&loginData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// get the token from service
		token, err := services.LoginUser(loginData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// return response if successfull
		ctx.JSON(http.StatusOK, token)
	}
}
