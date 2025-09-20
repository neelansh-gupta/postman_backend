package main

import (
	"log"
	"os"

	"go-api/handlers"
	"go-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get OMDb API key from environment
	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" {
		log.Fatal("OMDB_API_KEY environment variable is required")
	}

	// Initialize services
	omdbService := services.NewOMDbService(apiKey)

	// Initialize handlers
	movieHandler := handlers.NewMovieHandler(omdbService)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", movieHandler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Movie Details API - /api/movie?title=The Matrix
		api.GET("/movie", movieHandler.GetMovieDetails)

		// Episode Details API - /api/episode?series_title=Breaking Bad&season=1&episode_number=1
		api.GET("/episode", movieHandler.GetEpisodeDetails)

		// Genre-Based Movies API - /api/movies/genre?genre=Action
		api.GET("/movies/genre", movieHandler.GetMoviesByGenre)

		// Movie Recommendations API - /api/recommendations?favorite_movie=The Matrix
		api.GET("/recommendations", movieHandler.GetRecommendations)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET /health - Health check")
	log.Printf("  GET /api/movie?title=<movie_title> - Get movie details")
	log.Printf("  GET /api/episode?series_title=<series>&season=<num>&episode_number=<num> - Get episode details")
	log.Printf("  GET /api/movies/genre?genre=<genre> - Get top 15 movies by genre")
	log.Printf("  GET /api/recommendations?favorite_movie=<movie_title> - Get movie recommendations")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
