package middleware

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/models"

	jwt "github.com/golang-jwt/jwt/v5"
)

const authenticationDataKey = "middleware-auth-data"

func claimsFromToken(t *jwt.Token) (*models.Claims, error) {
	claims, ok := t.Claims.(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("failed get jwt.MapClaims from token.Claims")
	}

	return claims, nil
}

var (
	ErrNoAuthData        = fmt.Errorf("no auth data in current request context")
	ErrCorruptedAuthData = fmt.Errorf("corrupted auth data in current request context")
)

// TODO: fix nested if statements
func GetClaims(c *gin.Context) (*models.Claims, error) {
	if data, exists := c.Get(authenticationDataKey); !exists {
		return nil, ErrNoAuthData
	} else {
		authData, ok := data.(*models.Claims)
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

func NewAuthentication(receiverFunc jwt.Keyfunc) *AuthMiddleware {
	if receiverFunc == nil {
		panic("receiver cannot be nil")
	}

	return &AuthMiddleware{
		receiver: receiverFunc,
		log:      log.WithPrefix("authentication"),
	}
}

// TODO: Rename method
func (m *AuthMiddleware) Authentication(validators ...AuthValidateFunc) gin.HandlerFunc {
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

		claims, err := claimsFromToken(token)
		if err != nil {
			m.log.Error("Failed get claims from token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token content"})
			return
		}

		for _, v := range validators {
			if err := v(claims); err != nil {
				m.log.Error("Auth data validation failed, request aborted", "error", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not allowed"})
				return
			}
		}

		c.Set(authenticationDataKey, claims)

		c.Next()
	}
}

type AuthValidateFunc func(*models.Claims) error

func AdminOnly() AuthValidateFunc {
	return func(ad *models.Claims) error {

		// if !ad.IsAdmin() { return fmt.Errorf("current user is not an administrator, roles %v", ad.Roles)
		// }
		// return nil
		return nil
	}
}
