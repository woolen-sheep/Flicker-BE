package controller

import (
	"fmt"
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

	cardIDHex := c.Param("id")
	cardID, err := primitive.ObjectIDFromHex(cardIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	userIDHex := context.GetJWTUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return context.Error(c, http.StatusUnauthorized, "unauthorized", err)
	}

	cardsetIDHex := c.Param("cardset_id")

	if !isCardsetAccessible(cardsetIDHex, userIDHex) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	m := model.GetModel()
	defer m.Close()

	comment := model.Comment{
		OwnerID: userID,
		CardID:  cardID,
		Content: p.Comment,
		Status:  constant.StatusNormal,
	}
	commentIDHex, err := m.AddComment(comment)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when AddComment", err)
	}
	return context.Success(c, commentIDHex)
}

// GetComments will accept card ID and return comments of the certain card.
func GetComments(c echo.Context) error {
	cardIDHex := c.Param("id")
	cardsetIDHex := c.Param("cardset_id")
	userIDHex := context.GetJWTUserID(c)

	if !isCardsetAccessible(cardsetIDHex, userIDHex) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
	}

	m := model.GetModel()
	defer m.Close()

	// Test card ID
	_, cardExist, err := m.GetCard(cardIDHex)
	if !cardExist {
		return context.Error(c, http.StatusNotFound, "card not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCard", err)
	}

	comments, commentExist, err := m.GetComments(cardIDHex)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetComments", err)
	}
	if !commentExist {
		return context.Success(c, make(param.GetCommentResponse, 0))
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
			ID: cmt.ID.Hex(),
			Owner: param.UserResponse{
				ID:       user.ID.Hex(),
				Username: user.Username,
				Avatar:   user.Avatar,
			},
			Comment:    cmt.Content,
			LastUpdate: fmt.Sprintf("%d", cmt.LastUpdateTime),
		})
	}

	return context.Success(c, resp)
}

// DeleteComment will accept card ID & comment ID and delete the exact comment.
func DeleteComment(c echo.Context) error {
	cardIDHex := c.Param("id")
	cardID, err := primitive.ObjectIDFromHex(cardIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	commentIDHex := c.Param("comment_id")
	commentID, err := primitive.ObjectIDFromHex(commentIDHex)
	if err != nil {
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	userIDHex := context.GetJWTUserID(c)
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return context.Error(c, http.StatusUnauthorized, "unauthorized", err)
	}

	if !isCommentOwner(commentIDHex, userIDHex) {
		return context.Error(c, http.StatusForbidden, "permission denied", nil)
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

func isCommentOwner(comment, owner string) bool {
	m := model.GetModel()
	defer m.Close()
	_, err := m.GetCommentWithOwner(comment, owner)
	return err == nil
}
