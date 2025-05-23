package routes

import (
	"github.com/gin-gonic/gin"
	// Adjust import path to your anime-service module name
	"github.com/vrstep/wawatch-backend/controller"
	// Middleware specific to anime-service if needed
	// "github.com/vrstep/wawatch-anime-service/middleware"
)

// ProviderRoute defines routes for managing watch providers
func ProviderRoute(router *gin.Engine) {
	// Note: No RequireAuth middleware here by default.
	// You might add service-to-service auth later if needed.
	providers := router.Group("/providers")
	{
		// providers.POST("/", controller.AddWatchProvider)                  // Controller needs to be created/moved here
		providers.PUT("/:provider_id", controller.UpdateWatchProvider)    // Controller needs to be created/moved here
		providers.DELETE("/:provider_id", controller.DeleteWatchProvider) // Controller needs to be created/moved here
	}
}
