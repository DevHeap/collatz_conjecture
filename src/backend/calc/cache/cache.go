package cache

import (
	"math/big"
	"sync"
	"os"
	"log"
)

type Cache struct{
	sync.RWMutex
	maxCount  int // constant

	minLength int

	index   Index
	storage Storage

	db DB
}

const DefaultMaxCount = 300

func NewCacheDefault() *Cache {
	db, err := NewDatabase()

	if err != nil {
		log.Print("Failed to connect to to database: ", err)
		os.Exit(1)
	}

	return NewCache(DefaultMaxCount, db)
}

func NewCacheDummyDatabase() *Cache {
	db := NewDummyDatabase()
	return NewCache(DefaultMaxCount, db)
}

func NewCache(maxCount int, db DB) *Cache {
	cache := &Cache{
		maxCount:  maxCount,

		minLength: 0,

		index:   NewIndex(),
		storage: NewStorage(),

		db: db,
	}

	entries, err := db.Load()
	if err != nil {
		log.Print("Failed to load cache: ", err)
		os.Exit(1)
	}

	for _, value := range entries {
		cache.PutLocal(value)
	}


	log.Print("Cache initialized successfuly")
	return cache
}

func (c *Cache) Put(path []*big.Int) {
	c.put(path, true)
}

func (c *Cache) PutLocal(path []*big.Int) {
	c.put(path, false)
}

func (c *Cache) put(path []*big.Int, updateDb bool) {
	c.RLock()

	if len(path) < 1{
		// TODO log this is abnormal
		c.RUnlock()
		return
	}

	if !(c.index.Size() < c.maxCount || len(path) > c.minLength) {
		c.RUnlock()
		return
	}

	c.RUnlock()
	c.Lock()
	defer c.Unlock()

	key := path[0].String()

	// Insert into index
	c.index.Update(key, len(path))

	// If index contains more than max count, remove minimal
	if c.index.Size() > c.maxCount {
		minKey := c.index.MinKey()
		c.storage.Rem(minKey)
		c.index.RemMin()
		// Deletion should stay in case of race condition happened during
		// shutdown, resulting in having >300 entries in DB
		c.db.DeleteByKeyAsync(minKey)
	}

	// Putting new data into map
	c.storage.Put(path)

	// Updating minimal length
	c.minLength = c.index.MinLength()

	if updateDb {
		// Updating database
		c.db.InsertAsync(path[0], path)
	}
}

func (c *Cache) Get(number *big.Int) ([]*big.Int, bool) {
	c.RLock()
	defer c.RUnlock()

	return c.storage.Get(number)
}
