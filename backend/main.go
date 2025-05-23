package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/middleware"
	"github.com/vrstep/wawatch-backend/routes"
)

func main() {
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Logging()) // Enhanced logging middleware

	config.ConnectDB()

	// Apply CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// --- Route Setup ---
	routes.UserRoutes(router)             // User auth, profile routes
	routes.UserAnimeListRoutes(router)    // User's anime list management routes
	routes.AnimePassThroughRoutes(router) // Routes that proxy to anime-servic
	routes.UserHistoryRoutes(router)

	router.Run(":8080")
	// Remove the duplicate router.Run() line
}
