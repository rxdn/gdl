package cache

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/rxdn/gdl/objects/user"
	"os"
	"strconv"
	"sync"
)

type BoltCache struct {
	*bolt.DB
	Options CacheOptions

	// TODO: Should we store self in the DB? Seems kinda redundant
	selfLock sync.RWMutex
	self     user.User
}

type BoltOptions struct {
	ClearOnRestart bool
	Path string
	FileMode os.FileMode
	*bolt.Options
}

func NewBoltCache(cacheOptions CacheOptions, boltOptions BoltOptions) BoltCache {
	if boltOptions.ClearOnRestart {
		_ = os.Remove(boltOptions.Path)
	}

	db, err := bolt.Open(boltOptions.Path, boltOptions.FileMode, boltOptions.Options)
	if err != nil {
		panic(err)
	}

	if err := createBuckets(db); err != nil {
		panic(err)
	}

	return BoltCache{
		DB:      db,
		Options: cacheOptions,
	}
}

func createBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("users")); err != nil {
			return err
		}

		return nil
	})
}

func (c *BoltCache) StoreUser(u user.User) {
	c.StoreUsers([]user.User{u})
}

func (c *BoltCache) StoreUsers(users []user.User) {
	_ = c.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		for _, u := range users {
			if encoded, err := json.Marshal(u.ToCachedUser()); err == nil {
				if err := b.Put([]byte(strconv.FormatUint(u.Id, 10)), encoded); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})
}