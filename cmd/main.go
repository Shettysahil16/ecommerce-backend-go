package main

import (
	"backend/config"
	"backend/database"
	"backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 1. Connect to DB
	database.ConnectDB()

	// 2. Connect to Redis
	config.ConnectRedis()

	// 3. Create indexes
	database.CreateIndexes()

	// 4. Setup routes
	router := gin.Default()
	routes.AuthRoutes(router)

	// 5. Start server
	router.Run(":8080")
}
