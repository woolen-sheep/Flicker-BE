package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

func initCardsetRouter(g *echo.Group) {
	g.POST("", controller.NewCardset)
	g.PUT("/:id", controller.UpdateCardset)
	g.DELETE("/:id", controller.DeleteCardset)
	g.GET("/:id", controller.GetCardset)
}
