package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/middleware"
)

func UserHistoryRoutes(router *gin.Engine) {
	history := router.Group("/api/v1/me/history")
	history.Use(middleware.RequireAuth) // All history routes require authentication
	{
		history.GET("/", controller.GetUserViewHistory)
	}
}
