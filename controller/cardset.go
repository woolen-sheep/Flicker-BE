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
		return context.Error(c, http.StatusUnauthorized, "unauthorized", err)
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
		return context.Error(c, http.StatusUnauthorized, "unauthorized", err)
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
		return context.Error(c, http.StatusUnauthorized, "unauthorized", err)
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

	if cardset.Access != constant.CardsetAccessPublic && cardset.OwnerID.Hex() != userID {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	cardsIDs := []string{}
	cards, err := m.GetCardInCardset(cardsetID)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when get card id list", err)
	}
	for _, c := range cards {
		cardsIDs = append(cardsIDs, c.ID.Hex())
	}

	user, _, err := m.GetUser(cardset.OwnerID.Hex())
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when get cardset owner", err)
	}

	err = m.UpdateVisitCount(cardsetID)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when UpdateVisitCount", err)
	}

	resp := param.GetCardsetResponse{
		ID:            cardsetID,
		OwnerID:       cardset.OwnerID.Hex(),
		OwnerName:     user.Username,
		Name:          cardset.Name,
		Description:   cardset.Description,
		FavoriteCount: cardset.FavoriteCount,
		VisitCount:    cardset.VisitCount,
		Access:        cardset.Access,
		CreateTime:    cardset.CreateTime,
		Cards:         cardsIDs,
		IsFavorite:    m.IsUserFavorite(userID, cardset.ID),
	}
	return context.Success(c, resp)
}

// GetRandomCardsets will return a list of random public cardsets.
func GetRandomCardsets(c echo.Context) error {
	p := param.RandomCardsetsRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	err = c.Validate(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	cardsets, err := m.GetRandomCardset(p.Count)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCardset", err)
	}

	resp := []param.GetCardsetResponse{}
	for _, cs := range cardsets {
		resp = append(resp, param.GetCardsetResponse{
			ID:            cs.ID.Hex(),
			Name:          cs.Name,
			Description:   cs.Description,
			FavoriteCount: cs.FavoriteCount,
			VisitCount:    cs.VisitCount,
			CreateTime:    cs.CreateTime,
			Access:        cs.Access,
		})
	}
	return context.Success(c, resp)
}

// SearchCardsets will accept query param `keyword` and search in title and description of
// cardsets and return a cardset array.
func SearchCardsets(c echo.Context) error {
	p := param.SearchCardsetRequest{
		PageRequest: param.DefaultPageRequest,
	}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	cardsets, err := m.GetCardsetByKeyword(p.Keyword, p.Skip, p.Limit)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCardset", err)
	}

	resp := []param.GetCardsetResponse{}
	for _, cs := range cardsets {
		resp = append(resp, param.GetCardsetResponse{
			ID:            cs.ID.Hex(),
			Name:          cs.Name,
			Description:   cs.Description,
			FavoriteCount: cs.FavoriteCount,
			VisitCount:    cs.VisitCount,
			CreateTime:    cs.CreateTime,
			Access:        cs.Access,
		})
	}
	return context.Success(c, resp)
}

func isCardsetOwner(cardset, owner string) bool {
	m := model.GetModel()
	defer m.Close()
	_, err := m.GetCardsetWithOwner(cardset, owner)
	return err == nil
}

func isCardsetAccessible(cardset, user string) bool {
	m := model.GetModel()
	defer m.Close()

	_, err := m.GetCardsetWithOwner(cardset, user)
	if err == nil { // user is the owner of the cardset, then accessible
		return true
	}

	c, exist, err := m.GetCardset(cardset)
	if !exist || err != nil {
		return false
	} else { // cardset is public, then accessible
		return c.Access == 1
	}
}
