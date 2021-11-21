package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

func initCardRouter(g *echo.Group) {
	// Card APIs
	g.POST("", controller.NewCard)
	g.PUT("/:id", controller.UpdateCard)
	g.DELETE("/:id", controller.DeleteCard)
	g.GET("/:id", controller.GetCard)

	// Comment on card APIs
	g.POST("/:id/comment", controller.NewComment)
	g.GET("/:id/comment", controller.GetComments)
	g.DELETE("/:id/comment/:comment_id", controller.DeleteComment)
}
