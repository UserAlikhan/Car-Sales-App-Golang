package middlewares

import (
	"car_sales/internal/services"
	"car_sales/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckCarPostOwnershipMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ID from endpoint url
		ID, err := utils.GetIDParam(ctx, "ID", "carPostID")
		if err != nil || ID == 0 {
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

		ctx.Next()
	}
}

// Middleware that checks image ownership
func CheckImageOwnershipMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get ID parameter
		ID, err := utils.GetIDParam(ctx, "ID")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters."})
			ctx.Abort()
			return
		}

		// get params from authorization middleware
		isAdmin := ctx.GetBool("isAdmin")
		userID := ctx.GetInt("userID")

		// call image service to get image record
		carImage, err := services.GetCarImageByIDWithoutURL(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// only admin and image's owner could take action on image
		if !isAdmin && carImage.CarPost.SellerID != uint(userID) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized to take action with this record."})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
