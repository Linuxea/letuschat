package conn

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func NewChatMessageQueue() ChatMessageQueue {
	return &ChatMQ{
		redisCli: redis.NewClient(&redis.Options{
			Addr:         ":",
			MaxRetries:   3,
			DB:           0,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
			MaxConnAge:   time.Hour * 24 * 365,
		}),
	}
}

type ChatMessageQueue interface {
	Send([]byte) error
}

type ChatMQ struct {
	redisCli *redis.Client
}

func (mq *ChatMQ) Send(data []byte) error {
	return mq.redisCli.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: "chat",
		Values: data,
	}).Err()

}
