package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util"
	"github.com/woolen-sheep/Flicker-BE/util/context"
)

func SignUp(c echo.Context) error {
	p := param.SignUpRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	err = c.Validate(p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	code, err := m.GetVerifyCode(p.Mail)
	if err == model.ErrNotFound {
		return context.Error(c, http.StatusBadRequest, "code not found", err)
	}

	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "internal server error", err)
	}

	if p.Code != code {
		return context.Error(c, http.StatusBadRequest, "code error", err)
	}

	cipher, err := util.EncryptPassword(p.Password)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when encrypt password", err)
	}

	user := model.User{
		Mail:     p.Mail,
		Username: p.Username,
		Password: cipher,
	}
	_, err = m.AddUser(user)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when AddUser", err)
	}
	return context.Success(c, "ok")
}
