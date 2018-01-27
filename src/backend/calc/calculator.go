package calc

import (
	"math/big"
	"../calc/cache"
	"time"
	"sync"
)


type Calculator struct {
	DataCh chan Result

	next  *big.Int
	step  *big.Int
	guard sync.Mutex
	stopCh chan struct{}

	cache *cache.Cache
}

func NewCalculator(start *big.Int, workersCount int, cache *cache.Cache) Calculator {
	stopCh := make(chan struct{})
	dataCh := make(chan Result, 4)

	calculator := Calculator{
		DataCh: dataCh,

		next:  start,
		step:  new(big.Int).SetInt64(1),
		stopCh: stopCh,

		cache: cache,
	}

	for i := 0; i < workersCount; i++ {
		go calculator.compute(int64(i))
	}

	return calculator
}

func (c *Calculator) Stop() {
	close(c.stopCh)
}

func (c *Calculator) getNextNumber() *big.Int {
	c.guard.Lock()
	next := c.next
	c.next = c.next.Add(c.next, c.step)
	c.guard.Unlock()
	return next
}

func (c *Calculator) compute(offset int64) {
	for {
		select {
		default:
			number := c.getNextNumber()

			if path, ok := c.cache.Get(number); ok {
				c.DataCh <- NewResultFromPath(path, time.Duration(0))

			} else{
				start := time.Now()
				path  := FindPath(number)
				elapsed := time.Since(start)

				c.DataCh <- NewResultFromPath(path, elapsed)
				c.cache.Put(path)
			}
		case <-c.stopCh:
			return
		}
	}
}
