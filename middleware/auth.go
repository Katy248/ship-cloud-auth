package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/data"
)

var sessionKey = "session-" + uuid.New().String()

const AuthorizationHeader = "Authorization"

func WithAuthentication(ctx *gin.Context) {
	header := ctx.GetHeader(AuthorizationHeader)

	if header == "" {
		err := fmt.Errorf("unauthorized: header %q not specified", AuthorizationHeader)
		log.Error("No authorization header", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": err})
		return
	}

	if strings.Contains(header, "Bearer") {
		err := fmt.Errorf("token malformed: Armen includes 'Bearer' in tokent")
		log.Error("Армен, заебал, пиши авторизацию сам, а не ИИшкой", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": err, "armensMessage": "go fuck yourself"})
		return
	}

	token, err := jwt.ParseWithClaims(header, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return config.SecurityKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Error("JWT expired", "error", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "token expired"})
			return
		}
		log.Error("Failed parse JWT", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}
	if !token.Valid {
		log.Error("Invalid JWT", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}
	sessionIDstr := token.Claims.(*jwt.RegisteredClaims).ID

	sID, err := uuid.Parse(sessionIDstr)
	if err != nil {
		log.Error("Failed parse session ID", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}

	session, err := data.GetSession(sID)
	if err != nil {
		log.Error("Failed get session", "error", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "bad credentials"})
		return
	}

	if session.UserBlocked {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"details": "session blocked"})
	}

	ctx.Set(sessionKey, session)

}

func GetSession(ctx *gin.Context) *data.Session {
	session, ok := ctx.Get(sessionKey)
	if !ok {
		panic(fmt.Errorf("session not found (key %q), probably not authenticated", sessionKey))
	}
	return session.(*data.Session)
}
