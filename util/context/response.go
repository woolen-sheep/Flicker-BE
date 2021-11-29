package context

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/config"
	"github.com/woolen-sheep/Flicker-BE/constant/i18n"
)

// Response 返回值
type Response struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	ErrHint string      `json:"hint,omitempty"`
}

// Success 成功
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Data:    data,
		Error:   "",
		Success: true,
	})
}

// Error 错误
func Error(c echo.Context, status int, message string, err error) error {
	ret := Response{
		Data:    nil,
		Message: message,
		Success: false,
	}
	if config.C.I18N.Enabled {
		statusString, ok := i18n.Status[config.C.I18N.Language][status]
		if ok {
			ret.Message = statusString
		}
	}
	if err != nil {
		ret.Error = err.Error()
	}

	return c.JSON(status, ret)
}
