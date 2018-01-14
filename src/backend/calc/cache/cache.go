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


func NewCache() *Cache{
	return &Cache{
		maxCount:  300,

		minLength: 0,

		index:   NewIndex(),
		storage: NewStorage(),
	}
}

func (c *Cache) Put(number *big.Int, path []*big.Int) {
	c.Lock()
	defer c.Unlock()


	if !(c.index.Size() < c.maxCount || len(path) > c.minLength) {
		return
	}

	key := number.String()

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