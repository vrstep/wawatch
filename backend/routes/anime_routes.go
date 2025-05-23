package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/middleware"
)

// AnimePassThroughRoutes defines routes that forward requests to the anime-service.
// Most of these can be public if the anime-service itself doesn't require auth for them.
// Auth here is to protect the user-service endpoint itself if needed, not necessarily for the data.
func AnimePassThroughRoutes(router *gin.Engine) {
	// Group for all routes that are essentially proxies to the anime-service
	// The /api/v1 prefix matches what anime-service expects if client prepends it.
	// Or, make paths here identical to anime-service and client calls them directly.
	// Let's assume paths here mirror anime-service for clarity.
	proxiedAnime := router.Group("/ext/anime")
	proxiedAnime.Use(middleware.RequireAuth) // Using /ext to denote external call
	{
		proxiedAnime.GET("/search", controller.SearchAnime)
		proxiedAnime.GET("/popular", controller.GetPopularAnime)
		proxiedAnime.GET("/trending", controller.GetTrendingAnime)
		proxiedAnime.GET("/upcoming", controller.GetUpcomingAnime)                  // New in anime-service
		proxiedAnime.GET("/recently-released", controller.GetRecentlyReleasedAnime) // New
		proxiedAnime.GET("/explore", controller.ExploreAnime)                       // New
		proxiedAnime.GET("/season/:year/:season", controller.GetAnimeBySeason)

		// Recommendations might be user-specific eventually, but anime-service's is generic for now.
		// If it becomes personalized, anime-service would need user context (e.g. user ID).
		proxiedAnime.GET("/recommendations", controller.GetAnimeRecommendations)

		proxiedAnime.GET("/:id", controller.GetAnimeDetails) // Gets details & providers
	}
}
