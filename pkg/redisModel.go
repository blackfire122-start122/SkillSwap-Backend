package pkg

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx context.Context
var RedisClient *redis.Client

type RedisMessage struct {
	Message   string
	UserId    uint64
	ChatId    uint64
	CreatedAt time.Time `json:"created_at"`
}

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	Ctx = context.Background()
}
