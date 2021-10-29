package model

import (
	"context"
	"time"

	"github.com/woolen-sheep/Flicker-BE/config"
)

type ChatModel interface {
	AddMessage(msg, association, newerID string) error
	GetMessage(association, newerID string) ([]string, error)
	ClearMessage(association, newerID string) error
	IncrChattingCounter(association, newerID string) error
	DecrChattingCounter(association, newerID string) error
	GetChatCount(association, newerID string) (int, error)
	Close()
}

type Model interface {
	// Close will close database connection
	Close()
	// Abort will stop all statement and roll back when using transaction
	Abort()
	// TODO: add interfaces
}

type model struct {
	dbTrait
	ctx   context.Context
	abort bool
}

func GetModel() Model {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if config.C.Debug {
		ctx = context.Background()
	}

	ret := &model{
		dbTrait: getDBTx(ctx),
		ctx:     ctx,
		abort:   false,
	}

	return ret
}
