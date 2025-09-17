package handlers

import (
	"car_sales/internal/configs"
	"car_sales/internal/search"
	"car_sales/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchCarPostHandler(s3Conf *configs.S3Config) gin.HandlerFunc {
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

		// we are searching 2 times more images than asked, because in worse case we will have wholes in our data
		// we process the wholes and remove them at the service layer. So, if we removed one of them we will
		// get new image from the extra ones
		res, err := search.SearchCarPost(ctx.Request.Context(), q, from, size*2)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(res.Hits.Hits) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "No data found"})
			return
		}

		// after we got search result from elastic search, call service layer
		var carPostIDs []int
		for _, carPost := range res.Hits.Hits {
			carPostIDs = append(carPostIDs, int(carPost.Source.ID))
		}

		carPosts, err := services.GetCarPostsByIDsArray(ctx, s3Conf, carPostIDs, size)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, carPosts)
	}
}
