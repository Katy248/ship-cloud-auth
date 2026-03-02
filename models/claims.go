package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID  uuid.UUID `json:"userId"`
	Blocked bool      `json:"blocked"`
	Roles   []string  `json:"roles"`
}

// Deprecated: Do not use this method, use permissions handler
func (c *Claims) IsAdmin() bool {
	return false
}
