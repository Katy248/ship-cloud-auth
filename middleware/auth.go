package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/data"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/keyval"
)

var sessionKey = "session-" + uuid.New().String()

const AuthorizationHeader = "Authorization"

func WithAuthentication(ctx *gin.Context) {
	header := ctx.GetHeader(AuthorizationHeader)

	if header == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	token, err := jwt.ParseWithClaims(header, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Config.GetString("jwt-security-key"), nil
	})

	if err != nil {
		log.Error("Failed parse JWT", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}
	if !token.Valid {
		log.Error("Invalid JWT", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}
	sessionID := token.Claims.(jwt.RegisteredClaims).ID

	sessionJSON, err := keyval.RDB.Get(ctx.Request.Context(), sessionID).Result()
	if err != nil {
		log.Error("Failed receive session data from Redis", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}

	var session data.Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		log.Error("Failed unmarshal session data", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}

	if session.Blocked {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"details": "session blocked"})

	}

	ctx.Set(sessionKey, session)

}

func GetSession(ctx *gin.Context) data.Session {
	session, ok := ctx.Get(sessionKey)
	if !ok {
		panic(fmt.Errorf("session not found (key %q), probably not authenticated", sessionKey))
	}
	return session.(data.Session)
}
