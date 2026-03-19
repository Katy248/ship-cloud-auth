package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katy248/ship-cloud-auth/data"
)

func HandleGetPermissions(ctx *gin.Context) {
	permissions := data.GetAllPermissions()
	ctx.JSON(http.StatusOK, gin.H{"permissions": permissions})
}
