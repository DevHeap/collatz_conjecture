package calc

import (
	"math/big"
	"../calc/cache"
	"time"
)

type Calculator struct {
	DataCh chan Result

	start  *big.Int
	step   *big.Int
	stopCh chan struct{}

	cache *cache.Cache
}

func NewCalculator(start *big.Int, workersCount int, cache *cache.Cache) Calculator {
	stopCh := make(chan struct{})
	dataCh := make(chan Result, 4)

	calculator := Calculator{
		DataCh: dataCh,

		start:  start,
		step:   big.NewInt(int64(workersCount)),
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

func (c *Calculator) compute(offset int64) {
	number := new(big.Int).Set(c.start)
	number.Add(number, big.NewInt(offset))

	for {
		select {
		default:
			if path, ok := c.cache.Get(number); ok {
				c.DataCh <- NewResultFromPath(path, time.Duration(0))

			} else{
				start := time.Now()
				path  := FindPath(number)
				elapsed := time.Since(start)

				c.DataCh <- NewResultFromPath(path, elapsed)
				c.cache.Put(number, path)
			}

			number.Add(number, c.step)

		case <-c.stopCh:
			return
		}
	}
}
