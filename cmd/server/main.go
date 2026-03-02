package main

import (
	"fmt"
	"strings"

	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/internal/config"
	auth "sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/middleware"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/internal/database"

	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/internal/handlers"

	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/models"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Error occurred while loading .env file", "error", err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	viper.SetConfigName("cloud-auth-config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("Error occurred while reading config file:", err)
	}

}
func initDatabase(conf *config.Config) {
	db, err := database.InitDB(conf)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate models", "error", err)
	}
}

func main() {
	conf := config.New()
	initDatabase(conf)
	ensureAdminCreated(conf)

	authMiddleware := auth.NewAuthentication(conf.JWTSecret())

	server := gin.Default()

	{
		api := server.Group("/api/auth")
		api.POST("/login", handlers.Login)

		api.GET("/users", authMiddleware.Authentication(handlers.HasPermissions(models.ListUsersPermission)), handlers.ListUsers)
		api.GET("/users/me", authMiddleware.Authentication(), handlers.GetUser)
		api.GET("/users/:id", authMiddleware.Authentication(), handlers.GetUser)
		api.POST("/users/:id", authMiddleware.Authentication(), handlers.UpdateUser)
		api.POST("/users/:id/email", authMiddleware.Authentication(handlers.HasPermissions(models.SetEmailPermission)))
		api.POST("/users/:id/password", authMiddleware.Authentication(handlers.HasPermissions(models.SetEmailPermission)))
		api.POST("/users", handlers.CreateUser)
		api.POST("/users/:id/block", authMiddleware.Authentication(
			handlers.HasPermissions(models.BlockUserPermission)),
			handlers.BlockUser)
	}

	// Запуск сервера
	log.Info("Server starting", "port", conf.Port())
	if err := server.Run(fmt.Sprintf(":%d", conf.Port())); err != nil {
		log.Fatal("Failed to start server", "error", err)
	}
}

// ensureAdminCreated создает root пользователя если не существует
func ensureAdminCreated(conf *config.Config) {

	rootEmail := conf.RootUserEmail()
	rootPassword := conf.RootUserPassword()

	log.Info("Создание пользователя root", "email", rootEmail, "password", rootPassword)
	log.Warn("Измените данные пользователя root сразу же после запуска сервера")

	var adminRole models.Role
	if err := database.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Warn("Failed to get admin role, creating new", "error", err)
		adminRole = *createAdminRole()
	}

	var count int64
	database.DB.Model(&models.User{}).Where("email = ?", rootEmail).Count(&count)
	if count != 0 {
		log.Warn("User already exists, skipping", "email", rootEmail)
		return
	}
	user := models.User{
		Name:  "root",
		Email: rootEmail,
		Roles: models.Roles{adminRole},
	}

	if err := user.SetPassword(rootPassword); err != nil {
		log.Error("Failed to set password for root user", "error", err)
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Error("Failed to create root user", "error", err)
	} else {
		log.Info("Root user created successfully")
	}
}

func createAdminRole() *models.Role {
	role := models.Role{
		Name: "admin",
	}

	if err := database.DB.Create(&role).Error; err != nil {
		log.Fatal("Failed to create admin role", "error", err)
	} else {
		log.Info("Admin role created")
	}
	return &role
}
