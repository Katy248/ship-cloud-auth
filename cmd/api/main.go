package main

import (
	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/auth"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/data"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/db"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/handlers"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/keyval"
)

func main() {
	config.Setup()
	db.Setup()
	keyval.Setup()
	defer func() {
		if err := keyval.RDB.Close(); err != nil {
			log.Error("failed to close redis", "error", err)
		}
	}()

	data.Migrate()

	config.Config.RegisterAlias("jwt-security-key", "security-key")
	middleware := auth.DefaultMiddleware(config.Config)
	server := gin.Default()

	auth := server.Group("/api/auth")
	auth.POST("/register", handlers.HandleRegister)
	auth.POST("/login", handlers.HandleLogin)
	auth.POST("/refresh", middleware.WithMiddlewareOnly, handlers.HandleRefresh)

	users := server.Group("/api/users", middleware.WithAuthentication)
	users.GET("/:id", handlers.HandleGetUser)
	users.GET("/", handlers.HandleGetUsersList)
	users.POST("/:id/set-password", handlers.HandleUserSetPassword)
	users.POST("/:id/set-email", handlers.HandleUserSetEmail)
	users.POST("/:id/block", handlers.HandleUserBlock)

	roles := server.Group("/api/roles", middleware.WithAuthentication)
	roles.GET("/")
	roles.GET("/:id")

	server.GET("/api/permissions", handlers.HandleGetPermissions)

	if err := server.Run(":8080"); err != nil {
		log.Fatal("failed to start server", "error", err)
	}
}
