package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/woolen-sheep/Flicker-BE/config"
)

const (
	// jwt的过期时间，默认设置为7天
	jwtExpiresDuration = time.Hour * 24 * 7
)

// JWTClaims 使用的JWT结构，JWT的修改请直接修改结构中的字段
type JWTClaims struct {
	jwt.StandardClaims
	ID       string `json:"id"`
	Mail     string `json:"mail"`
	Username string `json:"username"`
}

// GenerateJWTToken 根据键值对生成jwt token
func GenerateJWTToken(claims JWTClaims) (string, error) {
	claims.ExpiresAt = time.Now().Add(jwtExpiresDuration).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(config.C.JWT.Secret))
}
