package middlewares

import (
	"car_sales/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckCarPostOwnershipMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ID from endpoint url
		ID, err := strconv.Atoi(ctx.Param("ID"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"errot": "Invalid parameters passed."})
			ctx.Abort()
			return
		}

		// get userID from authorization url
		isAdmin := ctx.GetBool("isAdmin")
		userID := ctx.GetInt("userID")

		// get car post by id
		carPost, err := services.GetCarPostByIDWithoutImageURLs(uint(ID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID. Car post was not found."})
			ctx.Abort()
			return
		}

		// only admin or the owner could proceed
		if !isAdmin && userID != int(carPost.SellerID) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to take action with this record."})
			ctx.Abort()
			return
		}
	}
}
