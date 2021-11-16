package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

func initUserRouter(g *echo.Group) {
	g.POST("", controller.SignUp)
	g.PUT("", controller.UpdateUser)
	g.GET("/:user_id", controller.GetUser)
	g.GET("/", controller.GetUser)
}
