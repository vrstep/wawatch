package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/middleware"
)

// AnimeRoute defines routes that forward requests to the anime-service
func AnimeRoute(router *gin.Engine) {
	anime := router.Group("/anime")
	{
		// These routes now point to controller functions that use AnimeClient
		// to call the anime-service. Authentication might still be required
		// here to protect the user service endpoint itself.
		anime.GET("/search", middleware.RequireAuth, controller.SearchAnime)
		anime.GET("/:id", controller.GetAnimeDetails) // Consider if this needs auth

		// Public discovery endpoints (forwarded)
		anime.GET("/popular", controller.GetPopularAnime)
		anime.GET("/trending", controller.GetTrendingAnime)
		anime.GET("/season/:year/:season", controller.GetAnimeBySeason)

		// Recommendations (forwarded)
		anime.GET("/recommendations", middleware.RequireAuth, controller.GetAnimeRecommendations)

		// Check if anime is in the *user's* list (handled by user-animelist controller)
		// This route might be better placed in user_animelist_routes.go, but if kept here:
		anime.GET("/:id/list-status", middleware.RequireAuth, controller.GetAnimeInUserList)

		// REMOVED: The route for adding providers is now handled by the anime-service
		// anime.POST("/provider", middleware.RequireAuth, controller.AddWatchProvider)
	}
}
