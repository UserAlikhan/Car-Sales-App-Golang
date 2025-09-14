package middlewares

import (
	"car_sales/internal/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	ipLimiterMap = sync.Map{}
)

// limit - number os requests allowed per second
// burst - maximum number of events that could happen at once
func RateLimiterMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch IP address
		ip := utils.GetIP(ctx.Request)

		// load or create limiter
		val, _ := ipLimiterMap.LoadOrStore(ip, rate.NewLimiter(limit, burst))
		limiter := val.(*rate.Limiter)

		// return error if the limit was reached
		if !limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please, try again later",
			})
			ctx.Abort()
			return
		}

		// if everything is good, forward the request
		ctx.Next()
	}
}
