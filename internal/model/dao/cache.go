package dao

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	Cache *bigcache.BigCache
)

func InitCache(ctx context.Context) *bigcache.BigCache {
	cache, err := bigcache.New(ctx, bigcache.DefaultConfig(30*time.Minute))
	if err != nil {
		panic(err)
	}
	Cache = cache
	return cache
}
