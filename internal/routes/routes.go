package routes

import (
	"car_sales/internal/handlers"
	"car_sales/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	api := r.Group("/carBrand")
	{
		api.GET("/getAllCarBrands", handlers.GetAllCarBrandsHandler())
		api.GET("/:id", handlers.GetCarBrandByIdHandler())
		api.POST("/createCarBrand", middlewares.AuthMiddleware(), middlewares.RequireAdminMiddleware(), handlers.CreateCarBrandHandler())
		api.PUT("/updateCarBrand", handlers.UpdateCarBrandHandler())
		api.DELETE("/deleteCarBrand", handlers.DeleteCarBrandHandler())
	}

	api = r.Group("/users")
	{
		api.POST("/signUp", handlers.SignUpHandler())
		api.POST("/login", handlers.LoginHandler())
	}
}
