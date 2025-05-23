package routes

import (
	"github.com/gin-gonic/gin"
	// Adjust import path to your anime-service module name
	"github.com/vrstep/wawatch-backend/controller"
	// Middleware specific to anime-service if needed (e.g., specific logging, rate limiting)
	// "github.com/vrstep/wawatch-anime-service/middleware"
)

// AnimeRoute defines routes related to fetching anime data
func AnimeRoute(router *gin.Engine) {
	// Note: No RequireAuth middleware here, as this service trusts the calling service (backend)
	anime := router.Group("/anime")
	{
		anime.GET("/search", controller.SearchAnime)  // Controller needs to be created/moved here
		anime.GET("/:id", controller.GetAnimeDetails) // Controller needs to be created/moved here

		// Public discovery endpoints
		anime.GET("/popular", controller.GetPopularAnime)               // Controller needs to be created/moved here
		anime.GET("/trending", controller.GetTrendingAnime)             // Controller needs to be created/moved here
		anime.GET("/season/:year/:season", controller.GetAnimeBySeason) // Controller needs to be created/moved here

		// Recommendations endpoint (implementation might differ from user service)
		anime.GET("/recommendations", controller.GetAnimeRecommendations)    // Controller needs to be created/moved here
		anime.GET("/upcoming", controller.GetUpcomingAnime)                  // Controller needs to be created/moved here
		anime.GET("/recently-released", controller.GetRecentlyReleasedAnime) // Controller needs to be created/moved here
		anime.GET("/explore", controller.ExploreAnime)                       // New explore endpoint

		anime.GET("/")
	}
}
