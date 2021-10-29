package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/woolen-sheep/Flicker-BE/config"
	"github.com/woolen-sheep/Flicker-BE/controller"
	middleware2 "github.com/woolen-sheep/Flicker-BE/middleware"
	"github.com/woolen-sheep/Flicker-BE/router"
	"github.com/woolen-sheep/Flicker-BE/util"
	. "github.com/woolen-sheep/Flicker-BE/util/log"
)

func main() {
	ok, err := util.ParseFlag()
	if err != nil {
		Logger.Fatal(err)
	}

	if !ok {
		return
	}

	e := echo.New()

	e.HTTPErrorHandler = controller.HTTPErrorHandler

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Validator = middleware2.GetValidator()
	err = middleware2.InitBeforeStart(e)
	if err != nil {
		Logger.Fatal(err)
	}

	gAPI := e.Group(config.C.App.Prefix)
	router.InitRouter(gAPI)

	Logger.Fatal(e.Start(config.C.App.Addr))
}
