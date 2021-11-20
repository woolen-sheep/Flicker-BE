package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

func initServiceRouter(g *echo.Group) {
	g.GET("/upload_token", controller.GetQiniuUploadToken)
}
