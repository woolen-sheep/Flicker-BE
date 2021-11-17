package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

// InitRouter will initialize all routers
func InitRouter(g *echo.Group) {
	initIndexRouter(g)
	usrGrp := g.Group("/user")
	initUserRouter(usrGrp)
	cardsetGrp := g.Group("/cardset")
	initCardsetRouter(cardsetGrp)
	cardGrp := g.Group("/cardset/:cardset_id/card")
	initCardRouter(cardGrp)
}

func initIndexRouter(g *echo.Group) {
	g.POST("/verify", controller.SendVerifyCode)
	g.POST("/login", controller.Login)
}
