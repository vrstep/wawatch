package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/middleware"
)

func UserAnimeListRoutes(router *gin.Engine) {
	// All user list operations require authentication
	list := router.Group("/api/v1/me/animelist") // User's own list
	list.Use(middleware.RequireAuth)
	{
		list.GET("/", controller.GetUserAnimeList)            // Get my list (paginated, status filter)
		list.POST("/", controller.AddToAnimeList)             // Add/Update anime in my list
		list.PATCH("/entry/:id", controller.UpdateListEntry)  // Update specific fields of a list entry (by list entry DB ID)
		list.DELETE("/entry/:id", controller.DeleteListEntry) // Delete a list entry (by list entry DB ID)
		list.GET("/stats", controller.GetUserAnimeListStats)

		// Check status of a specific anime (by its external ID) in the user's list
		list.GET("/status/:animeExternalID", controller.GetAnimeInUserList)
	}
}
