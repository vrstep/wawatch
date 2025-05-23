package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/middleware"
)

func UserRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/signup", controller.Signup)
		auth.POST("/login", controller.Login)
		auth.GET("/validate", middleware.RequireAuth, controller.Validate)
		// router.POST("/logout", controller.Logout) // You'd need a Logout handler
	}

	profile := router.Group("/api/v1/me/profile")
	profile.Use(middleware.RequireAuth)
	{
		profile.GET("/", controller.GetMyProfile)
		profile.PUT("/", controller.UpdateMyProfile)
		profile.PUT("/password", controller.ChangeMyPassword)
		profile.PUT("/username", controller.ChangeMyUsername)
		profile.PUT("/email", controller.ChangeMyEmail)
	}

	// Public user views
	usersPublic := router.Group("/api/v1/users")
	{
		// usersPublic.GET("/:username/profile", controller.GetUserPublicProfile) // New
		usersPublic.GET("/:username/animelist", controller.GetUserPublicAnimeList)
	}

	// Admin routes for user management - protect these with an admin role middleware
	// adminUsers := router.Group("/api/v1/admin/users")
	// adminUsers.Use(middleware.RequireAuth, middleware.RequireAdminRole)
	// {
	//    adminUsers.GET("/", controller.GetUsers)
	//    adminUsers.POST("/", controller.CreateUser) // Usually signup is public
	//    adminUsers.GET("/:id", controller.GetUser)
	//    adminUsers.PUT("/:id", controller.UpdateUser)
	//    adminUsers.DELETE("/:id", controller.DeleteUser)
	// }
}
