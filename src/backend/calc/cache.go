package calc

import (
	"encoding/json"
	"math/big"
	"sort"
	"sync"
	"log"
)

type record struct {
	key    string
	length int
}

type Cache struct {
	sync.RWMutex
	maxCount  int
	minLength int

	index     []record
	storage   map[string][]byte
}

func NewCache() *Cache {

	return &Cache{
		maxCount: 300,
		minLength: 0,

		index:   make([]record, 0),
		storage: make(map[string][]byte),
	}
}

func (c *Cache) Get(number *big.Int) (Result, bool) {
	key := number.String()

	c.RLock()
	defer c.RUnlock()

	if byt, ok := c.storage[key]; ok {
		log.Print("Cache hit")
		var result Result
		if err := json.Unmarshal(byt, &result); err != nil {
			panic(err)
		}
		return result, true
	}
	return Result{}, false
}

func (c *Cache) Put(number *big.Int, result *Result) {
	c.Lock()
	defer c.Unlock()

	// Check length, call for real put only if needed
	if result.PathLength > c.minLength {
		c.put(number, result)
		log.Print("Cache update")
	}
}

func (c *Cache) put(number *big.Int, result *Result) {
	key := number.String()
	data, err := json.Marshal(result)

	if err != nil {
		// this should not happen
		panic(err)
	}

	// insert into index
	c.updateIndex(record{key: key, length: result.PathLength})

	// if index contains more than max count, remove minimal
	if len(c.index) > c.maxCount {
		// removing minimal element from map
		delete(c.storage, c.index[0].key)

		// yes, this should delete minimal element in index
		c.index = append(c.index[:0], c.index[1:]...)
	}

	// putting new data into map
	c.storage[key] = data

	// updating minimal length
	c.minLength = c.index[0].length

	// TODO update database here
}

// Helper function to insert element into index and maintain order in it
func (c *Cache) updateIndex(element record) {
	// search for proper position of element in index
	position := sort.Search(len(c.index), func(i int) bool { return c.index[i].length > element.length })
	// extend index
	c.index = append(c.index, record{})
	// shift elements by copy
	copy(c.index[position+1:], c.index[position:])
	// put new element on proper position
	c.index[position] = element
}
