package main

import (
	"github.com/gin-gonic/gin"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/data"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/db"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/handlers"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/keyval"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/middleware"
)

func main() {
	config.Setup()
	db.Setup()
	keyval.Setup()
	defer keyval.RDB.Close()

	data.Migrate()

	server := gin.Default()

	users := server.Group("/api/users", middleware.WithAuthentication)

	auth := server.Group("/api/auth")
	auth.POST("/register", handlers.HandleRegister)
	auth.POST("/login")
	auth.POST("/refresh")

	users.GET("/:id", handlers.HandleGetUser)
	users.GET("/", handlers.HandleGetUsersList)
	users.POST("/:id/set-password", handlers.HandleUserSetPassword)
	users.POST("/:id/set-email", handlers.HandleUserSetEmail)
	users.POST("/:id/block", handlers.HandleUserBlock)

	roles := server.Group("/api/roles")

	roles.GET("/")
	roles.GET("/:id")

	sessions := server.Group("/api/sessions")
	sessions.GET("/current", handlers.HandleGetSession)
	sessions.GET("/", handlers.HandleGetSessionsList)
	sessions.DELETE("/:id", handlers.HandleDeleteSession)

	server.Run(":8080")
}
