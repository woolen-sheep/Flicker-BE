package model

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	QiniuUploadRecordDuration = time.Minute * 60 //七牛上传凭证的默认时长为1小时
)

type QiniuUpload struct {
	User primitive.ObjectID `json:"user"`
	Url  string             `json:"url"`
}

func AddQiniuUploadRecord(key string, upload QiniuUpload) error {
	jsonStr, err := json.Marshal(upload)
	if err != nil {
		return err
	}
	result := redisClient.Set(key, jsonStr, QiniuUploadRecordDuration)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
