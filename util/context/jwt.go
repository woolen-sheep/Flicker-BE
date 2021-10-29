package context

import (
	"github.com/dgrijalva/jwt-go"
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
	return getJWTFiled(c, "user_id")
}

func getJWTFiled(c echo.Context, filedName string) string {
	token := c.Get(config.JWTContextKey)
	if token != nil {
		if tokenStr, ok := token.(*jwt.Token).Claims.(jwt.MapClaims)[filedName].(string); ok {
			return tokenStr
		}
		return ""
	}
	return ""
}
