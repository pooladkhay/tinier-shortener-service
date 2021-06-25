package client

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pooladkhay/tinier-shortener-service/domain"
)

var redisOpts domain.RedisOpts

func init() {
	// env.Load()
	ctx := context.Background()

	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RDB_ADDR"),
		Password: os.Getenv("RDB_PASSWORD"),
		DB:       0,
	})
	redisOpts.Client = r
	redisOpts.Ctx = &ctx

}

func NewRedis() *domain.RedisOpts {
	return &redisOpts
}
