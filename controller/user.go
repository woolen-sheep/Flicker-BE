package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util"
	"github.com/woolen-sheep/Flicker-BE/util/context"
)

// SignUp will check mail verify code and add a new user.
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

// UpdateUser will update user info of current user. Empty fields will be ignored.
func UpdateUser(c echo.Context) error {
	p := param.UpdateUserRequest{}
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

	fmt.Println(context.GetJWTUserID(c))
	userID := context.GetJWTUserID(c)

	oldUser, exist, err := m.GetUser(userID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "user not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetUser", err)
	}

	var cipher string
	if len(p.Password) != 0 {
		cipher, err = util.EncryptPassword(p.Password)
		if err != nil {
			return context.Error(c, http.StatusInternalServerError, "error when encrypt password", err)
		}
	} else {
		cipher = ""
	}

	user := model.User{
		ID:       oldUser.ID,
		Username: p.Username,
		Password: cipher,
	}
	err = m.UpdateUser(user)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when UpdateUser", err)
	}
	return context.Success(c, "ok")
}

// GetUser will accept user ID and return user info. If param `user_id`
// is empty, will try to use user ID in JWT.
func GetUser(c echo.Context) error {
	userID := c.Param("user_id")

	if len(userID) == 0 {
		userID = context.GetJWTUserID(c)
	}

	m := model.GetModel()
	defer m.Close()

	user, exist, err := m.GetUser(userID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "user not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetUser", err)
	}
	resp := param.UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Avatar:   user.Avatar,
	}
	for _, f := range user.Favorite {
		resp.Favorite = append(resp.Favorite, f.Hex())
	}
	return context.Success(c, resp)
}

// UpdateFavorite will add a cardset into current user collection
func UpdateFavorite(c echo.Context) error {
	p := param.AddCollectionRequest{}
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

	userID := context.GetJWTUserID(c)

	err = m.UpdateFavorite(userID, p.CardsetID, p.Liked)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when UpdateUser", err)
	}
	return context.Success(c, "ok")
}

// GetFavorite will return favorite cardsets of current user.
func GetFavorite(c echo.Context) error {
	userID := c.Param("user_id")

	if len(userID) == 0 {
		userID = context.GetJWTUserID(c)
	}

	m := model.GetModel()
	defer m.Close()

	user, exist, err := m.GetUser(userID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "user not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetUser", err)
	}

	cardsets, err := m.GetCardsetByIDList(user.Favorite)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCardsetByIDList", err)
	}
	res := []param.CardsetInfoResponse{}
	for _, cs := range cardsets {
		res = append(res, param.CardsetInfoResponse{
			ID:          cs.ID.Hex(),
			OwnerID:     cs.OwnerID.Hex(),
			Name:        cs.Name,
			Description: cs.Description,
			Access:      cs.Access,
		})
	}
	return context.Success(c, res)
}

// GetCreated will return cardsets created by current user.
func GetCreated(c echo.Context) error {
	userID := c.Param("user_id")

	if len(userID) == 0 {
		userID = context.GetJWTUserID(c)
	}

	m := model.GetModel()
	defer m.Close()

	user, exist, err := m.GetUser(userID)
	if !exist {
		return context.Error(c, http.StatusNotFound, "user not found", nil)
	}
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetUser", err)
	}

	cardsets, err := m.GetCardsetByIDList(user.Favorite)
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when GetCardsetByIDList", err)
	}
	res := []param.CardsetInfoResponse{}
	for _, cs := range cardsets {
		res = append(res, param.CardsetInfoResponse{
			ID:          cs.ID.Hex(),
			OwnerID:     cs.OwnerID.Hex(),
			Name:        cs.Name,
			Description: cs.Description,
			Access:      cs.Access,
		})
	}
	return context.Success(c, res)
}
