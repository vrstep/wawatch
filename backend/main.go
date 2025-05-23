package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/middleware"
	"github.com/vrstep/wawatch-backend/routes"
)

func main() {
	router := gin.New()
	router.Use(middleware.Logging())
	config.ConnectDB()

	router.Use(func(c *gin.Context) {
		// Determine the frontend origin. The error message indicates 'http://localhost'.
		// This could be http://localhost (port 80) or a specific port like http://localhost:5173.
		// Adjust `envFrontendOrigin` or its default based on your actual frontend setup.
		// The error "from origin 'http://localhost'" suggests the browser sees 'http://localhost' as the origin.
		defaultFrontendOrigin := "http://localhost"
		envFrontendOrigin := os.Getenv("CORS_FRONTEND_ORIGIN")
		if envFrontendOrigin == "" {
			envFrontendOrigin = defaultFrontendOrigin
			log.Printf("Warning: CORS_FRONTEND_ORIGIN not set, defaulting to %s. Verify this matches your frontend's actual origin.", envFrontendOrigin)
		}

		requestOrigin := c.Request.Header.Get("Origin")
		allowedOrigin := ""

		if requestOrigin == envFrontendOrigin {
			allowedOrigin = envFrontendOrigin
		} else if os.Getenv("GIN_MODE") != "release" && requestOrigin != "" && (strings.HasPrefix(requestOrigin, "http://localhost:") || requestOrigin == "http://localhost") {
			// For development, allow other localhost origins (e.g., http://localhost:3000, http://localhost:5173)
			log.Printf("CORS [DEV]: Allowing specific origin %s for request", requestOrigin)
			allowedOrigin = requestOrigin
		}

		// Only set Access-Control-Allow-Origin if we have a valid, non-wildcard origin
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			// Crucially, set Allow-Credentials to true only when a specific origin is allowed
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if requestOrigin != "" {
			// If an origin was provided in the request but it's not in our allowed list.
			log.Printf("CORS: Request from origin '%s' is not explicitly allowed. Configured frontend origin is '%s'.", requestOrigin, envFrontendOrigin)
			// Do not set Access-Control-Allow-Origin to '*' if you expect credentials.
			// The browser will block if it's a credentialed request and origin doesn't match.
		}
		// IMPORTANT: Remove or comment out the unconditional wildcard setting:
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // THIS WAS THE PROBLEM for credentialed requests

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			// Preflight requests should also have the correct CORS headers set above.
			// If allowedOrigin is not set, and it's a credentialed preflight, it will be correctly blocked by the browser.
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	routes.UserRoutes(router)
	routes.UserAnimeListRoutes(router)
	routes.AnimePassThroughRoutes(router)
	routes.UserHistoryRoutes(router)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
