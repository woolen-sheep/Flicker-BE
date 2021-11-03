package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util"
	"github.com/woolen-sheep/Flicker-BE/util/context"
)

func SendVerifyCode(c echo.Context) error {
	p := param.VerifyCodeRequest{}

	if err := c.Bind(&p); err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	if err := c.Validate(p); err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	blocking, err := m.VerifyCodeBlocking(p.Mail)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "internal server error", err)
	}

	if blocking {
		return context.Error(c, http.StatusBadRequest, "code send too often", nil)
	}

	code := util.GenerateVerifyCode()
	err = m.SetVerifyCode(p.Mail, code)
	if err != nil {
		m.Abort()
		return context.Error(c, http.StatusInternalServerError, "error when set code", err)
	}

	err = util.SendMailVerifyCode(p.Mail, code)
	if err != nil {
		m.Abort()
		return context.Error(c, http.StatusInternalServerError, "error when send mail", err)
	}

	return context.Success(c, "ok")
}

func Login(c echo.Context) error {
	p := param.LoginRequest{}
	if err := c.Bind(&p); err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	if err := c.Validate(p); err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	user, err := m.GetUserByMail(p.Mail)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "error when get user", err)
	}

	ok, err := util.ComparePassword(user.Password, p.Password)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when compare", err)
	}
	if !ok {
		return context.Error(c, http.StatusBadRequest, "wrong password", nil)
	}

	claim := util.JWTClaims{
		ID:       user.ID.Hex(),
		Mail:     user.Mail,
		Username: user.Username,
	}

	token, err := util.GenerateJWTToken(claim)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "internal server error", err)
	}

	return context.Success(c, token)
}
