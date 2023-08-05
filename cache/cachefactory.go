package cache

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type CacheFactory func() Cache

func PgCacheFactory(db *pgxpool.Pool, options CacheOptions) CacheFactory {
	return func() Cache {
		c := NewPgCache(db, options)
		return &c
	}
}

func MemoryCacheFactory(cacheOptions CacheOptions) CacheFactory {
	return func() Cache {
		c := NewMemoryCache(cacheOptions)
		return &c
	}
}
