package search

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchCarPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		q := ctx.Query("q")
		if q == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No query provided."})
			return
		}

		// Default Query takes "page" parameter from query,
		// but if it is not provided sets 1
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

		if page < 1 {
			page = 1
		}

		if size < 1 {
			size = 10
		}

		from := (page - 1) * size

		res, err := SearchCarPost(ctx.Request.Context(), q, from, size)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}
