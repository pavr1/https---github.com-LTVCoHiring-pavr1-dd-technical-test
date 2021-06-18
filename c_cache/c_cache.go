package c_cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CCache struct {
	Cache *cache.Cache
}

func New() *CCache {
	cache := &CCache{
		Cache: cache.New(720*time.Hour, 60*time.Minute),
	}

	return cache
}

func (c *CCache) Set(key string, data interface{}) {
	c.Cache.Set(key, data, cache.NoExpiration)
}

func (c *CCache) Get(key string) (interface{}, bool) {
	return c.Cache.Get(key)
}

func (c *CCache) Exist(key string) bool {
	_, found := c.Cache.Get(key)

	return found
}
