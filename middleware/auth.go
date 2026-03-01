package middleware

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/models"

	jwt "github.com/golang-jwt/jwt/v5"
)

const AuthDataKey = "middleware-auth-data"

var (
	ErrNoAuthData        = fmt.Errorf("no auth data in current request context")
	ErrCorruptedAuthData = fmt.Errorf("corrupted auth data in current request context")
)

// TODO: fix nested if statements
func GetUser(c *gin.Context) (*AuthData, error) {
	if data, exists := c.Get(AuthDataKey); !exists {
		return nil, ErrNoAuthData
	} else {
		authData, ok := data.(*AuthData)
		if !ok {
			return nil, ErrCorruptedAuthData
		} else {
			return authData, nil
		}
	}
}

type AuthMiddleware struct {
	receiver jwt.Keyfunc
	log      *log.Logger
}

func New(receiverFunc jwt.Keyfunc) *AuthMiddleware {
	if receiverFunc == nil {
		panic("receiver cannot be nil")
	}

	return &AuthMiddleware{
		receiver: receiverFunc,
		log:      log.WithPrefix("authentication"),
	}
}

// TODO: Rename method
func (m *AuthMiddleware) Check(validators ...AuthValidateFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			m.log.Error("Authentication failed with no token specified in Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token specified"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, m.receiver)
		if err != nil {
			m.log.Error("Invalid token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		data, err := NewDataFromToken(token)
		if err != nil {
			m.log.Error("Failed get authentication data from token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token content"})
			return
		}

		for _, v := range validators {
			if err := v(*data); err != nil {
				m.log.Error("Auth data validation failed, request aborted", "error", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
				return
			}
		}

		c.Set(AuthDataKey, data)

		c.Next()
	}
}

type AuthValidateFunc func(AuthData) error

func AdminOnly() AuthValidateFunc {
	return func(ad AuthData) error {
		if !ad.IsAdmin() {
			return fmt.Errorf("current user is not an administrator, roles %v", ad.Roles)
		}
		return nil
	}
}
