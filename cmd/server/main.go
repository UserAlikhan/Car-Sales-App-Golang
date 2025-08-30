package main

import (
	"car_sales/internal/configs"
	"car_sales/internal/database"
	"car_sales/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	configs.LoadEnvVariables()
	database.ConnectToDB()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is not set or empty. Using the default value.")
		port = "8080"
	}

	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r)

	r.Run(":" + port)
}
