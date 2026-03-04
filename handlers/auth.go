package handlers

import (
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/data"
)

func HandleRegister(c *gin.Context) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	user, err := data.NewUser(
		request.Name, request.Email, request.Password,
	)
	if err != nil {
		if errors.Is(err, data.EmailAlreadyTakenErr) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"details": "email already taken",
			})
			return
		}
		log.Error("Failed create new user", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}
