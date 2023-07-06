package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var DBCache *cache.Cache

func init() {
	initDbCache()
}

func initDbCache() *cache.Cache {
	c := cache.New(24*time.Hour, 24*time.Hour)
	DBCache = c

	return DBCache
}
