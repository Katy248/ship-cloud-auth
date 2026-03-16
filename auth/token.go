package auth

import (
	"errors"
	"fmt"

	"charm.land/log/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Returns error if token is invalid or expired
func (m *Middleware) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, m.keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Error("JWT expired", "error", err)
			return nil, fmt.Errorf("token expired")
		}
		log.Error("Failed parse JWT", "error", err)
		return nil, fmt.Errorf("failed parse token: %s", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}
	return claims, nil
}
