package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

func initRecordRouter(g *echo.Group) {
	// Card APIs
	g.POST("/:cardset_id/:card_id", controller.UpdateRecord)
	g.GET("/:cardset_id", controller.GetRecords)
	g.DELETE("/:cardset_id", controller.ClearRecords)
}
