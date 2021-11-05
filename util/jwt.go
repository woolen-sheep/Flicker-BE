package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/woolen-sheep/Flicker-BE/config"
)

const (
	jwtExpiresDuration = time.Hour * 24 * 7
)

// JWTClaims with custom fields
type JWTClaims struct {
	jwt.StandardClaims
	ID       string `json:"id"`
	Mail     string `json:"mail"`
	Username string `json:"username"`
}

// GenerateJWTToken by claims
func GenerateJWTToken(claims JWTClaims) (string, error) {
	claims.ExpiresAt = time.Now().Add(jwtExpiresDuration).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(config.C.JWT.Secret))
}
