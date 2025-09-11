package main

import (
	"car_sales/internal/configs"
	"car_sales/internal/database"
	"car_sales/internal/routes"
	"car_sales/internal/search"
	"context"
	"fmt"
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

	// Redis cache initialization
	configs.InitRedis()

	// Setup s3 uploader
	s3Conf, err := configs.NewS3Config(context.TODO(), os.Getenv("BUCKET_NAME"), os.Getenv("AWS_REGION"))
	if err != nil {
		fmt.Printf("Failed to initialize S3 config %v", err)
		return
	}

	// initialize elastic search
	search.InitElasticSearch()

	// ensure index exists
	if err := search.CreateIndex(context.Background()); err != nil {
		fmt.Println("Error creating elastic search index: ", err)
	}

	// Initialize routes
	routes.InitRoutes(r, s3Conf)

	r.Run(":" + port)
}
