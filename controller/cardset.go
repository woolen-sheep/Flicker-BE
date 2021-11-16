package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/constant"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	userID, err := primitive.ObjectIDFromHex(context.GetJWTUserID(c))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	cardset := model.Cardset{
		OwnerID:     userID,
		Name:        p.Name,
		Description: p.Description,
		Access:      p.Access,
		Status:      constant.StatusNormal,
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

	cardsetID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	userID, err := primitive.ObjectIDFromHex(context.GetJWTUserID(c))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	cardset := model.Cardset{
		OwnerID:     userID,
		ID:          cardsetID,
		Name:        p.Name,
		Description: p.Description,
		Access:      p.Access,
	}

	err = m.UpdateCardset(cardset)
	if err == model.ErrNotFound {
		return context.Error(c, http.StatusNotFound, "cardset not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when UpdateCardset", err)
	}
	return context.Success(c, "ok")
}

// DeleteCardset will accept cardset ID and delete the exact cardset.
func DeleteCardset(c echo.Context) error {
	cardsetID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	userID, err := primitive.ObjectIDFromHex(context.GetJWTUserID(c))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	cardset := model.Cardset{
		ID:      cardsetID,
		OwnerID: userID,
		Status:  constant.StatusDeleted,
	}

	exist, err := m.DeleteCardset(cardset)
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

	userID := context.GetJWTUserID(c)

	var cardsIDs []string
	if cardset.Access != constant.CardsetAccessPublic && cardset.OwnerID.Hex() != userID {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	resp := param.GetCardsetResponse{
		ID:          cardsetID,
		Name:        cardset.Name,
		Description: cardset.Description,
		Access:      cardset.Access,
		Cards:       cardsIDs,
	}
	return context.Success(c, resp)
}

func isCardsetOwner(cardset, owner string) bool {
	m := model.GetModel()
	defer m.Close()
	_, err := m.GetCardsetWithOwner(cardset, owner)
	return err == nil
}
