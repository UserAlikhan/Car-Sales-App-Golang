package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileCheckerMiddleware(maxSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// gets all the sent MultipartForm files
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form"})
			return
		}

		// parsing each file
		for _, files := range form.File {
			for _, fileHeader := range files {
				// check the size of the file
				if fileHeader.Size > maxSize {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf(
							"File %s is too large. Maximum allowed size is %d MD",
							fileHeader.Filename, maxSize,
						),
					})
					ctx.Abort()
					return
				}
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Files were uploaded successfully."})
		// proceed
		ctx.Next()
	}
}
