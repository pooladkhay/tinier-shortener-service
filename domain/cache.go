package domain

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisOpts struct {
	Client *redis.Client
	Ctx    *context.Context
}
