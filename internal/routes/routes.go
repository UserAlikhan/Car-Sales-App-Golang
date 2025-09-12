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
		api.GET("/getCarModels", handlers.GetAllCarModelsHandler())
		api.GET("getCarModels/:carBrandID", handlers.GetCarModelsByCarBrandIDHandler())
		api.POST("/createCarModel", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.CreateCarModelHandler())
		api.PUT("/updateCarModel/:ID", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.UpdateCarModelHandler())
		api.DELETE("deleteCarModel/:ID", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.DeleteCarModelHandler())
	}

	// Route and supporting endpoints for /carPost
	api = r.Group("/carPost")
	{
		api.POST(
			"/createCarPost",
			middlewares.AuthMiddleware(),
			// 7<<20 means 7MB in bytes
			middlewares.FileSizeCheckerMiddleware(7<<20),
			middlewares.ImageFileExtensionChecker(),
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
		api.PUT("/updateCarPost/:ID", handlers.UpdateCarPostHandler())
	}

	// Route and supporting endpoints for /carPostsImages
	api = r.Group("/carPostsImages")
	{
		// this endpoint allows to upload images to the existing car post
		// only authorized owners or admin can add images into someone's car post
		api.POST(
			"/uploadImages/:carPostID",
			middlewares.AuthMiddleware(),
			middlewares.CheckCarPostOwnershipMiddleware(),
			middlewares.FileSizeCheckerMiddleware(7<<20),
			middlewares.ImageFileExtensionChecker(),
			handlers.UploadImagesHandler(s3Conf),
		)
		// only authorized owners or admin can delete someone's images
		api.DELETE(
			"/deleteSingleImage/:ID",
			middlewares.AuthMiddleware(),
			middlewares.CheckImageOwnershipMiddleware(),
			handlers.DeleteCarPostImageHandler(s3Conf),
		)
	}
}
