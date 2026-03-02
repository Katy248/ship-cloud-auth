package server

import (
	"github.com/gin-gonic/gin"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/internal/config"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/internal/handlers"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/middleware"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/models"
)

func New(conf *config.Config) *gin.Engine {

	authMiddleware := middleware.NewAuthentication(conf.JWTSecret())
	server := gin.Default()

	{
		api := server.Group("/api/auth")
		api.POST("/login", handlers.Login)

		api.GET("/users", authMiddleware.Authentication(handlers.HasPermissions(models.ListUsersPermission)), handlers.ListUsers)
		api.GET("/users/me", authMiddleware.Authentication(), handlers.GetUser)
		api.GET("/users/:id", authMiddleware.Authentication(), handlers.GetUser)
		api.POST("/users/:id", authMiddleware.Authentication(), handlers.UpdateUser)
		api.POST("/users/:id/email", authMiddleware.Authentication(handlers.HasPermissions(models.SetEmailPermission)))
		api.POST("/users/:id/password", authMiddleware.Authentication(handlers.HasPermissions(models.SetPasswordPermission)))
		api.POST("/users", handlers.CreateUser)
		api.POST("/users/:id/block", authMiddleware.Authentication(
			handlers.HasPermissions(models.BlockUserPermission)),
			handlers.BlockUser)
	}

	return server
}
