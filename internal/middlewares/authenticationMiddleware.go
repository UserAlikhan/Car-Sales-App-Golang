package middlewares

import (
	"car_sales/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token and check it
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Missing or Invalid Authorization header"})
			return
		}
		// split main part of the token from Bearer prefix
		token := strings.TrimPrefix(authorizationHeader, "Bearer ")
		if token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Missing Authorization header"})
			ctx.Abort()
			return
		}

		// validate token and get claims back
		claims, err := utils.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}
		// set variables received from decoded token
		ctx.Set("userID", claims.ID)
		ctx.Set("isAdmin", claims.IsAdmin)
		// proceed
		ctx.Next()
	}
}
