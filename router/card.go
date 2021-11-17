package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

func initCardRouter(g *echo.Group) {
	g.POST("", controller.NewCard)
	g.PUT("/:id", controller.UpdateCard)
	g.DELETE("/:id", controller.DeleteCard)
	g.GET("/:id", controller.GetCard)
}
