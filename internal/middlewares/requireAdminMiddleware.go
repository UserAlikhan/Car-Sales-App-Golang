package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// this middleware goes after AuthMiddleware for routes
// that need bo accessed only by admins
func RequireAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin, exists := ctx.Get("isAdmin")
		if !exists || isAdmin != true {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
