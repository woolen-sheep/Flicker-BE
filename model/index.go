package model

import (
	"context"
	"time"

	"github.com/woolen-sheep/Flicker-BE/config"
)

type objArray []interface{}

type Model interface {
	// Close will close database connection
	Close()
	// Abort will stop all statement and roll back when using transaction
	Abort()
	// VerifyCode
	VerifyCodeBlocking(mail string) (bool, error)
	SetVerifyCode(mail, code string) error
	GetVerifyCode(mail string) (string, error)
	// UserInterface contains all user functions in model layer
	UserInterface
	// CardsetInterface contains all cardset functions in model layer
	CardsetInterface
	// CardInterface contains all card functions in model layer
	CardInterface
	// CommentInterface contains all comment functions in model layer
	CommentInterface
	// RecordInterface contains all record functions in model layer
	RecordInterface
}

type model struct {
	dbTrait
	ctx    context.Context
	abort  bool
	cancel context.CancelFunc
}

func GetModel() Model {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	if config.C.Debug {
		ctx = context.Background()
	}

	ret := &model{
		dbTrait: getDBTx(ctx),
		ctx:     ctx,
		abort:   false,
		cancel:  cancel,
	}

	return ret
}
