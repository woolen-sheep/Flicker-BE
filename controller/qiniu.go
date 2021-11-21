package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"github.com/woolen-sheep/Flicker-BE/config"
	"github.com/woolen-sheep/Flicker-BE/controller/param"
	"github.com/woolen-sheep/Flicker-BE/model"
	"github.com/woolen-sheep/Flicker-BE/util"
	"github.com/woolen-sheep/Flicker-BE/util/context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetQiniuUploadToken(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(context.GetJWTUserID(c))
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when get user id", err)
	}

	typ := c.QueryParam("type")

	key := ""
	id := uuid.NewV4().String()

	switch typ {
	case "avatar":
		key = fmt.Sprintf("avatar/%s", id)
	case "image":
		key = fmt.Sprintf("image/%s", id)
	case "audio":
		key = fmt.Sprintf("audio/%s", id)
	default:
		return context.Error(c, http.StatusBadRequest, "type error", err)
	}

	resp := param.GetQiniuTokenRequest{
		URL:         fmt.Sprintf("%s/%s", config.C.Qiniu.URL, key),
		ResourceKey: key,
		Token:       util.GetQiniuOverwriteUpToken(key),
	}

	err = model.AddQiniuUploadRecord(resp.ResourceKey, model.QiniuUpload{
		User: userID,
		Url:  resp.URL,
	})
	if err != nil {
		return context.Error(c, http.StatusInternalServerError, "error when add qiniu record", err)
	}

	return context.Success(c, resp)
}
