package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
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
		if errors.Is(err, data.ErrEmailAlreadyTaken) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"details": "email already taken",
			})
			return
		}
		log.Error("Failed create new user", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func HandleLogin(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	user, err := data.GetUserByEmail(request.Email)
	if err != nil {
		log.Error("Failed get user by email", "error", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "invalid credentials"})
		return
	}

	if !user.ComparePassword(request.Password) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"details": "invalid credentials",
		})
		return
	}

	session, err := data.NewSession(user.ID)
	if err != nil {
		log.Error("Failed create new session", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	log.Info("New session", "session", session)

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"token":        createJWT(session),
		"refreshToken": createRefreshJWT(session),
	})
}

func HandleRefresh(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}
	session, err := data.GetSession(uuid.MustParse(request.RefreshToken))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"details": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        createJWT(session),
		"refreshToken": createRefreshJWT(session),
	})

}

const tokenTTL = time.Minute * 10
const refreshTokenTTL = time.Hour * 24

func newJWT(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString(config.SecurityKey())
	if err != nil {
		panic(fmt.Errorf("failed sign JWT: %s", err))
	}
	return signed
}

func createJWT(session *data.Session) string {

	claims := jwt.RegisteredClaims{
		ID:        session.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
	}
	return newJWT(claims)
}
func createRefreshJWT(session *data.Session) string {
	claims := jwt.RegisteredClaims{
		ID:        session.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
	}

	return newJWT(claims)
}
