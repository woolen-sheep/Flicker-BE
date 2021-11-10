package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util/context"
)

// NewCardset will add a new cardset.
func NewCardset(c echo.Context) error {
	p := param.NewCardsetRequest{}
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

	cardset := model.Cardset{
		Name:        p.Name,
		Description: p.Description,
		Access:      p.Access,
	}
	cardsetID, err := m.AddCardset(cardset)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when AddCardset", err)
	}
	return context.Success(c, cardsetID)
}

// UpdateCardset will update cardset info of current cardset. Empty fields will be ignored.
func UpdateCardset(c echo.Context) error {
	p := param.UpdateCardsetRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	oldCardset, exist, err := m.GetCardset(c.Param("id"))
	if !exist {
		return context.Error(c, http.StatusNotFound, "cardset not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCardset", err)
	}

	cardset := model.Cardset{
		ID:          oldCardset.ID,
		Name:        p.Name,
		Description: p.Description,
		Access:      p.Access,
	}

	err = m.UpdateCardset(cardset)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when UpdateCardset", err)
	}
	return context.Success(c, "ok")
}

// DeleteCardset will accept cardset ID and delete the exact cardset.
func DeleteCardset(c echo.Context) error {
	cardsetID := c.Param("id")

	m := model.GetModel()
	defer m.Close()

	exist, err := m.DeleteCardset(cardsetID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "cardset not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when DeleteCardset", err)
	}

	return context.Success(c, "ok")
}

// GetCardset will accept cardset ID and return cardset info.
func GetCardset(c echo.Context) error {
	cardsetID := c.Param("id")

	m := model.GetModel()
	defer m.Close()

	cardset, exist, err := m.GetCardset(cardsetID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "cardset not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCardset", err)
	}
	resp := param.GetCardsetResponse{
		ID:          cardsetID,
		Name:        cardset.Name,
		Description: cardset.Description,
		Access:      cardset.Access,
		Cards:       cardset.Cards,
	}
	return context.Success(c, resp)
}
