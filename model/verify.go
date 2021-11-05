package model

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/woolen-sheep/Flicker-BE/config"
)

func configPrefix(p string) string {
	return fmt.Sprintf("VerifyCode:%s", p)
}

func blockPrefix(p string) string {
	return fmt.Sprintf("Blocking:%s", p)
}

// VerifyCodeBlocking will check the duration after last mail
// and return true if the duration is too short.
func (m *model) VerifyCodeBlocking(mail string) (bool, error) {
	n, err := redisClient.Exists(blockPrefix(mail)).Result()
	return n != 0, err
}

// SetVerifyCode sets verify code of mail.
func (m *model) SetVerifyCode(mail, code string) error {
	_, err := redisClient.Set(configPrefix(mail), code, config.DurationCodeExpire).Result()
	if err != nil {
		return err
	}
	_, err = redisClient.Set(blockPrefix(mail), code, config.DurationCodeResend).Result()
	return err
}

// GetVerifyCode returns verify code of `mail` and returns `ErrNotFound` when
// there is no code of mail.
func (m *model) GetVerifyCode(mail string) (string, error) {
	code, err := redisClient.Get(configPrefix(mail)).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}

	return code, err
}
