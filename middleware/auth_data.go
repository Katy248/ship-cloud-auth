package middleware

import (
	"fmt"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/models"
)

type AuthData struct {
	UserID uuid.UUID
	Roles  []string
}

const AdminRole = "admin"

func (a *AuthData) IsAdmin() bool {
	return slices.Contains(a.Roles, AdminRole)
}

func NewDataFromToken(t *jwt.Token) (*AuthData, error) {
	claims, ok := t.Claims.(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("failed get jwt.MapClaims from token.Claims")
	}

	authData := &AuthData{
		UserID: claims.UserID,
		Roles:  claims.Roles,
	}

	return authData, nil
}
