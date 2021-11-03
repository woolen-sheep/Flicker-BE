package context

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/config"
	"github.com/woolen-sheep/Flicker-BE/util"
)

// GetJWTClaims returns JWTClaims
func GetJWTClaims(c echo.Context) *util.JWTClaims {
	return c.Get(config.JWTContextKey).(*jwt.Token).Claims.(*util.JWTClaims)
}

// GetJWTUserID returns only user_id
func GetJWTUserID(c echo.Context) string {
	return GetJWTClaims(c).ID
}
