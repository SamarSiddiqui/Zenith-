package main

import (
	"fmt"
	"log"
	"os"

	"habit-tracker/internal/db"
	"habit-tracker/internal/handlers/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Habit Tracker Backend...")

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, proceeding with environment variables")
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	if err := db.ConnectMongoDB(mongoURI); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Disconnect()

	// Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Habit Tracker Backend is healthy!",
		})
	})

	// Setup Routes
	user.SetupUserRoutes(r)

	// Start server
	log.Printf("Listening on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
