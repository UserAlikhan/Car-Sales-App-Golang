package routes

import (
	"car_sales/internal/configs"
	"car_sales/internal/handlers"
	"car_sales/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, s3Conf *configs.S3Config) {
	// Route and supporting endpoints for /carBrand
	api := r.Group("/carBrand")
	{
		api.GET("/getAllCarBrands", handlers.GetAllCarBrandsHandler())
		api.GET("/:id", handlers.GetCarBrandByIdHandler())
		// user must be authorized and must be an admin
		api.POST("/createCarBrand", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.CreateCarBrandHandler())
		api.PUT("/updateCarBrand/:id", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.UpdateCarBrandHandler())
		api.DELETE(
			"/deleteCarBrand/:id",
			middlewares.AuthMiddleware(),
			middlewares.RequireAdminMiddleware(),
			handlers.DeleteCarBrandHandler(),
		)
		api.POST(
			"/createCarBrandWithModels",
			middlewares.AuthMiddleware(),
			middlewares.RequireAdminMiddleware(),
			handlers.CreateCarBrandWithModelsHandler(),
		)
		api.POST(
			"/uploadLogo/:id",
			middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(),
			middlewares.FileSizeCheckerMiddleware(7<<20),
			middlewares.ImageFileExtensionChecker(),
			handlers.UploadLogoHandler(s3Conf),
		)
	}

	// Route and supporting endpoints for /users
	api = r.Group("/users")
	{
		api.POST("/signUp", handlers.SignUpHandler())
		api.POST("/login", handlers.LoginHandler())
	}

	// Route and supporting endpoints for /carModel
	api = r.Group("/carModel")
	{
		api.POST("/createCarModel", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.CreateCarModelHandler())
	}

	// Route and supporting endpoints for /carPost
	api = r.Group("/carPost")
	{
		// 7<<20 means 7MB in bytes
		api.POST(
			"/createCarPost",
			middlewares.AuthMiddleware(),
			middlewares.FileSizeCheckerMiddleware(7<<20),
			middlewares.ImageFileExtensionChecker(),
			// middlewares.CheckCarPostOwnershipMiddleware(),
			handlers.CreateCarPostHandler(s3Conf),
		)
		api.GET("/getAllUsersCarPosts/:userId", handlers.GetAllUsersCarPostsHandler())
		// only authorized owner or admin car delete someone's car post
		api.DELETE(
			"/deleteCarPost/:ID",
			middlewares.AuthMiddleware(),
			middlewares.CheckCarPostOwnershipMiddleware(),
			handlers.DeleteCarPostHandler(s3Conf),
		)
		api.GET("/getCarPostByID/:ID", handlers.GetCarPostByIDHandler(s3Conf))
		api.GET("/getCarPosts", handlers.GetCarPostsWithPaginationHandler(s3Conf))
	}

	// Route and supporting endpoints for /carPostsImages
	api = r.Group("/carPostsImages")
	{
		api.POST("/uploadImages", middlewares.ImageFileExtensionChecker(), handlers.UploadImagesHandler())
		// only authorized owners or admin can delete someone's images
		api.DELETE("/deleteSingleImage", middlewares.AuthMiddleware(), handlers.DeleteCarPostImageHandler(s3Conf))
	}
}
