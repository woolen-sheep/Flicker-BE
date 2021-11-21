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

// NewComment will append a comment to a certain card.
func NewComment(c echo.Context) error {
	p := param.NewCommentRequest{}
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

	cardID, err := primitive.ObjectIDFromHex(c.Param("id"))

	m := model.GetModel()
	defer m.Close()

	comment := model.Comment{
		OwnerID: userID,
		CardID:  cardID,
		Content: p.Comment,
		Status:  constant.StatusNormal,
	}
	commentID, err := m.AddComment(comment)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when AddComment", err)
	}
	return context.Success(c, commentID)
}

// GetComments will accept card ID and return comments of the certain card.
func GetComments(c echo.Context) error {
	cardID := c.Param("id")

	m := model.GetModel()
	defer m.Close()

	// Test card ID
	_, cardExist, err := m.GetCard(cardID)
	if !cardExist {
		return context.Error(c, http.StatusNotFound, "card not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCard", err)
	}

	comments, commentExist, err := m.GetComments(cardID)
	if !commentExist {
		// something to log
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetComments", err)
	}

	var resp param.GetCommentResponse
	for _, cmt := range comments {
		user, exist, err := m.GetUser(cmt.OwnerID.Hex())
		if !exist {
			return context.Error(c, http.StatusNotFound, "user not found", nil)
		}
		if err != nil {
			return context.Error(c, http.StatusInternalServerError, "error when GetUser", err)
		}

		resp = append(resp, param.CommentResponseItem{
			Owner: param.UserResponse{
				ID:       user.ID.Hex(),
				Username: user.Username,
				Avatar:   user.Avatar,
			},
			Comment: cmt.Content,
		})
	}

	return context.Success(c, resp)
}

// DeleteComment will accept card ID & comment ID and delete the exact comment.
func DeleteComment(c echo.Context) error {
	cardID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}
	commentID, err := primitive.ObjectIDFromHex(c.Param("comment_id"))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	userID, err := primitive.ObjectIDFromHex(context.GetJWTUserID(c))
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	m := model.GetModel()
	defer m.Close()

	comment := model.Comment{
		ID:      commentID,
		OwnerID: userID,
		CardID:  cardID,
		Status:  constant.StatusDeleted,
	}

	exist, err := m.DeleteComment(comment)
	if !exist {
		return context.Error(c, http.StatusNotFound, "comment not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when DeleteComment", err)
	}

	return context.Success(c, "ok")
}
