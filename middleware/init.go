// Package middleware contains middleware echo needed
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/woolen-sheep/Flicker-BE/config"
)

// InitBeforeStart 会在Web服务启动之前对echo实例进行一些初始化操作
func InitBeforeStart(e *echo.Echo) error {
	// 使用JWT
	e.Use(middleware.JWTWithConfig(CustomJWTConfig(config.C.JWT.Skip, "Bearer")))
	// 使用cors
	e.Use(middleware.CORS())
	return nil
}
