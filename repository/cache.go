package repository

import (
	"time"

	"github.com/pooladkhay/tinier-shortener-service/domain"
	"github.com/pooladkhay/tinier-shortener-service/helper/errs"
)

type Cache interface {
	Cache(hash, url string, exp int)
	GetCache(hash string) (*string, *errs.Err)
	Delete(hash string)
}

type cache struct {
	rdb *domain.RedisOpts
}

func NewCache(r *domain.RedisOpts) Cache {
	return &cache{rdb: r}
}

func (c *cache) Cache(hash, url string, exp int) {
	c.rdb.Client.Set(*c.rdb.Ctx, hash, url, time.Second*time.Duration(exp)).Err()
}

func (c *cache) GetCache(hash string) (*string, *errs.Err) {
	result, err := c.rdb.Client.Get(*c.rdb.Ctx, hash).Result()
	if err != nil {
		return nil, nil
	}
	return &result, nil
}

func (c *cache) Delete(hash string) {
	c.rdb.Client.Del(*c.rdb.Ctx, hash).Err()
}
