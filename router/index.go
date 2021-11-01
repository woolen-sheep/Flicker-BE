package router

import (
	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller"
)

// InitRouter will initialize all routers
func InitRouter(g *echo.Group) {
	initIndexRouter(g)
	grp := g.Group("/user")
	initUserRouter(grp)
}

func initIndexRouter(g *echo.Group) {
	g.POST("/verify", controller.SendVerifyCode)
}
