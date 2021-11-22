package controller

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/constant"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewCard will add a new card.
func NewCard(c echo.Context) error {
	p := param.NewCardRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	err = c.Validate(p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	cardsetIDHex := c.Param("cardset_id")
	if len(cardsetIDHex) == 0 {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	cardsetID, err := primitive.ObjectIDFromHex(cardsetIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	card := model.Card{
		CardsetID: cardsetID,
		Question:  p.Question,
		Answer:    p.Answer,
		Image:     p.Image,
		Audio:     p.Audio,
		Status:    constant.StatusNormal,
	}
	cardID, err := m.AddCard(card)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when AddCard", err)
	}
	return context.Success(c, cardID)
}

// UpdateCard will update card info of current card. Empty fields will be ignored.
func UpdateCard(c echo.Context) error {
	cardsetID := c.Param("cardset_id")

	userID := context.GetJWTUserID(c)

	if !isCardsetOwner(cardsetID, userID) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}
	p := param.UpdateCardRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	oldCard, exist, err := m.GetCard(c.Param("id"))
	if !exist {
		return context.Error(c, http.StatusNotFound, "card not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCard", err)
	}

	card := model.Card{
		ID:       oldCard.ID,
		Question: p.Question,
		Answer:   p.Answer,
		Image:    p.Image,
		Audio:    p.Audio,
	}

	err = m.UpdateCard(card)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when UpdateCard", err)
	}
	return context.Success(c, "ok")
}

// DeleteCard will accept card ID and delete the exact card.
func DeleteCard(c echo.Context) error {
	cardIDHex := c.Param("id")
	cardsetIDHex := c.Param("cardset_id")

	userID := context.GetJWTUserID(c)

	if !isCardsetOwner(cardsetIDHex, userID) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	m := model.GetModel()
	defer m.Close()

	cardID, err := primitive.ObjectIDFromHex(cardIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	cardsetID, err := primitive.ObjectIDFromHex(cardsetIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	card := model.Card{
		ID:        cardID,
		CardsetID: cardsetID,
		Status:    constant.StatusDeleted,
	}

	exist, err := m.DeleteCard(card)
	if !exist {
		return context.Error(c, http.StatusNotFound, "card not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when DeleteCard", err)
	}

	return context.Success(c, "ok")
}

// GetCard will accept card ID and return card info.
func GetCard(c echo.Context) error {
	cardID := c.Param("id")
	cardsetID := c.Param("cardset_id")

	userID := context.GetJWTUserID(c)

	if !isCardsetOwner(cardsetID, userID) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	m := model.GetModel()
	defer m.Close()

	card, exist, err := m.GetCard(cardID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "card not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCard", err)
	}

	resp := param.GetCardResponse{
		ID:       card.ID.Hex(),
		Question: card.Question,
		Answer:   card.Answer,
		Image:    card.Image,
		Audio:    card.Audio,
	}
	return context.Success(c, resp)
}

// GetCards will accept a card ID list and return a card info list.
func GetCards(c echo.Context) error {
	cardIDHexes := []string{}
	err := json.Unmarshal([]byte(c.QueryParam("ids")), &cardIDHexes)
	if err != nil || len(cardIDHexes) == 0 {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	cardsetID := c.Param("cardset_id")

	userID := context.GetJWTUserID(c)

	if !isCardsetOwner(cardsetID, userID) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	m := model.GetModel()
	defer m.Close()
	cardIDs := []primitive.ObjectID{}
	for _, hex := range cardIDHexes {
		id, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			return context.Error(c, http.StatusBadRequest, "bad request", err)
		}
		cardIDs = append(cardIDs, id)
	}
	cards, err := m.GetCardByIDList(cardIDs)
	if err == model.ErrNotFound {
		return context.Error(c, http.StatusNotFound, "card not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCard", err)
	}
	resp := []param.GetCardResponse{}
	for _, card := range cards {
		resp = append(resp, param.GetCardResponse{
			ID:       card.ID.Hex(),
			Question: card.Question,
			Answer:   card.Answer,
			Image:    card.Image,
			Audio:    card.Audio,
		})
	}

	return context.Success(c, resp)
}
