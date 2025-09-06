package middlewares

import (
	"car_sales/internal/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware that checks file size based on specified size
func FileSizeCheckerMiddleware(maxSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// gets all the sent MultipartForm files
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form."})
			ctx.Abort()
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

		// if everything is good, proceed
		ctx.Next()
	}
}

// Middleware that checks if image file has a valid extension
func ImageFileExtensionChecker() gin.HandlerFunc {
	// allowed file extensions array
	allowedImageFileExtensions := []string{"png", "jpg", "jpeg"}

	return func(ctx *gin.Context) {
		// get sent MultipartForm files data
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form."})
			ctx.Abort()
			return
		}

		// iterate through all files
		for _, files := range form.File {
			for _, fileHeader := range files {
				// get a filename
				filename := fileHeader.Filename
				// split extension part from the filename
				parts := strings.Split(filename, ".")
				if len(parts) < 2 {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "File has no extension."})
					ctx.Abort()
					return
				}

				// get extension part and convert to lower
				extension := strings.ToLower(parts[len(parts)-1])

				// custom Find function that checks if a variable is inside of an array
				isValid := utils.Find(allowedImageFileExtensions, extension)
				if !isValid {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf(
							"%s does not have allowed extension. Your file must have either .png .jpg or .jpeg",
							filename,
						),
					})
					ctx.Abort()
					return
				}
			}
		}

		// if everything is good, proceed
		ctx.Next()
	}
}
