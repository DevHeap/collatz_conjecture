package cache

import (
	"math/big"
	"sync"
)

type Cache struct{
	sync.RWMutex
	maxCount  int // constant

	minLength int

	index   Index
	storage Storage
}

const DefaultMaxCount = 300

func NewCacheDefault() *Cache {
	return NewCache(DefaultMaxCount)
}

func NewCache(maxCount int) *Cache{
	return &Cache{
		maxCount:  maxCount,

		minLength: 0,

		index:   NewIndex(),
		storage: NewStorage(),
	}
}

func (c *Cache) Put(path []*big.Int) {
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
		c.storage.Rem(c.index.MinKey())
		c.index.RemMin()
	}

	// Putting new data into map
	c.storage.Put(path)

	// Updating minimal length
	 c.minLength = c.index.MinLength()

	// TODO update database here
}



func (c *Cache) Get(number *big.Int) ([]*big.Int, bool) {
	c.RLock()
	defer c.RUnlock()

	return c.storage.Get(number)
}