package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	jwt "github.com/golang-jwt/jwt/v5"
)

const AuthDataKey = "middleware-auth-data"

const EnvironmentJWTKey = "JWT_SECURITY_KEY"

func EnvironmentReceiver() jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		key := os.Getenv(EnvironmentJWTKey)
		if key == "" {
			return nil, fmt.Errorf("environment variable %q not specified", EnvironmentJWTKey)
		}
		return []byte(key), nil
	}
}

func getJwtKey() jwt.Keyfunc {

	return func(t *jwt.Token) (any, error) {
		key := os.Getenv(EnvironmentJWTKey)
		if key == "" {
			return nil, fmt.Errorf("environment variable %q not specified", EnvironmentJWTKey)
		}
		return []byte(key), nil
	}
}

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

func (m *AuthMiddleware) Check(validators ...AuthValidateFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			m.log.Error("Authentication failed with no token specified in Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token specified"})
			return
		}
		token, err := jwt.Parse(tokenString, getJwtKey())
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

func Middleware(validators ...AuthValidateFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			log.Error("Authentication failed with no token specified in Authorization header")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, getJwtKey())
		if err != nil {
			log.Error("Authentication failed", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		data, err := NewDataFromToken(token)
		if err != nil {
			log.Error("Failed get authentication data from token", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		for _, v := range validators {
			if err := v(*data); err != nil {
				log.Error("Auth data validation failed, request aborted", "error", err)
				c.AbortWithStatus(http.StatusUnauthorized)
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
