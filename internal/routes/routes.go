package routes

import (
	"car_sales/internal/handlers"
	"car_sales/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	// Routes and supporting endpoints for /carBrand
	api := r.Group("/carBrand")
	{
		api.GET("/getAllCarBrands", handlers.GetAllCarBrandsHandler())
		api.GET("/:id", handlers.GetCarBrandByIdHandler())
		// user must be authorized and must be an admin
		api.POST("/createCarBrand", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.CreateCarBrandHandler())
		api.PUT("/updateCarBrand/:id", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.UpdateCarBrandHandler())
		api.DELETE("/deleteCarBrand/:id", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.DeleteCarBrandHandler())
		api.POST("/createCarBrandWithModels", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.CreateCarBrandWithModelsHandler())
	}

	// Routes and supporting endpoints for /users
	api = r.Group("/users")
	{
		api.POST("/signUp", handlers.SignUpHandler())
		api.POST("/login", handlers.LoginHandler())
	}

	// Routes and supporting endpoints for /carModel
	api = r.Group("/carModel")
	{
		api.POST("/createCarModel", handlers.CreateCarModel())
	}

}
