package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/util/context"
)

// HTTPErrorHandler will handle all errors and convert errors to certain format
func HTTPErrorHandler(err error, c echo.Context) {
	httpError, ok := err.(*echo.HTTPError)
	if ok {
		_ = context.Error(c, httpError.Code, fmt.Sprintf("%s", httpError.Message), err)
		return
	}

	_ = context.Error(c, http.StatusInternalServerError, "Unhandled internal server error", err)
}
