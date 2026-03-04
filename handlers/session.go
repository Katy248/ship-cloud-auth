package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/data"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/middleware"
)

func HandleGetSession(c *gin.Context) {
	session := middleware.GetSession(c)

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func HandleGetSessionsList(c *gin.Context) {

	session := middleware.GetSession(c)

	sessions, err := data.GetActiveSessions(session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}
func HandleDeleteSession(c *gin.Context) {
	sessionIDstr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDstr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	session, err := data.GetSession(sessionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	if session.UserID != middleware.GetSession(c).UserID {
		log.Error("User try to delete not own session",
			"user.id", session.UserID,
			"session.id", sessionID,
			"session.user.id", session.UserID)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	err = data.DeleteSession(sessionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
