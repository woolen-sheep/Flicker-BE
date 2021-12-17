package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateRecord will accept card ID and update study record
func UpdateRecord(c echo.Context) error {
	cardIDHex := c.Param("card_id")
	cardsetIDHex := c.Param("cardset_id")

	userIDHex := context.GetJWTUserID(c)

	p := param.UpdateRecordRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return context.Error(c, http.StatusUnauthorized, "wrong user id", err)
	}
	cardID, err := primitive.ObjectIDFromHex(cardIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	cardsetID, err := primitive.ObjectIDFromHex(cardsetIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	record := model.Record{
		OwnerID:   userID,
		CardsetID: cardsetID,
		CardID:    cardID,
		Status:    p.Status,
	}

	err = m.UpdateRecord(record)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "update record failed", err)
	}
	return context.Success(c, "ok")
}

// GetRecords will accept cardset ID and return all study record of current user.
func GetRecords(c echo.Context) error {
	cardsetIDHex := c.Param("cardset_id")

	userIDHex := context.GetJWTUserID(c)

	p := param.UpdateRecordRequest{}
	err := c.Bind(&p)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return context.Error(c, http.StatusUnauthorized, "wrong user id", err)
	}
	cardsetID, err := primitive.ObjectIDFromHex(cardsetIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	resp := param.GetRecordsRespons{Records: []param.RecordResponse{}}
	resp.Total, err = m.GetCount(cardsetIDHex)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "count cards failed", err)
	}

	records, err := m.GetRecords(cardsetID, userID)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "get records failed", err)
	}

	for _, r := range records {
		resp.Records = append(resp.Records, param.RecordResponse{
			CardID:     r.CardID.Hex(),
			LastStudy:  r.LastStudy,
			StudyTimes: r.StudyTimes,
			Status:     r.Status,
		})
	}
	return context.Success(c, resp)
}
