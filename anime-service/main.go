package main

import (
	"github.com/gin-gonic/gin"
	// Adjust these import paths based on your actual module name for anime-service
	"github.com/vrstep/wawatch-backend/config"     // Anime service's config
	"github.com/vrstep/wawatch-backend/middleware" // Anime service's middleware
	"github.com/vrstep/wawatch-backend/routes"     // Anime service's routes
)

func main() {
	// Use gin.New() for more control over middleware
	router := gin.New()

	// --- Middleware Setup ---

	// 1. RequestID Middleware: Ensures every request has an ID (generates if missing)
	router.Use(middleware.RequestID()) // Assuming you create this in anime-service/middleware

	// 2. Logging Middleware: Logs request details including RequestID, Time, Duration
	router.Use(middleware.Logging()) // Use the enhanced logging middleware

	// 3. CORS Middleware: Allow requests from your frontend
	router.Use(func(c *gin.Context) {
		// Replace "*" with your frontend origin in production for security
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// Ensure X-Request-ID is allowed and exposed
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH") // Added PATCH
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")                          // Expose the Request ID

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// --- Database Connection & Migrations ---
	// Connect to the database specific to the anime service
	// This ConnectDB should also handle running migrations for anime_caches, watch_providers
	config.ConnectDB()

	// --- Route Setup ---
	// Register routes handled by this service
	routes.AnimeRoute(router)    // Routes like /anime/search, /anime/:id, /anime/popular etc.
	routes.ProviderRoute(router) // Routes like /providers/:id (PUT, DELETE)

	// --- Start Server ---
	// Run on a different port than the main backend service
	router.Run(":8082")
}
