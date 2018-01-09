package calc

import (
	"math/big"
)

type Calculator struct {
	DataCh chan Result

	start  *big.Int
	step   *big.Int
	stopCh chan struct{}

	cache *Cache
}

func NewCalculator(start *big.Int, workersCount int) Calculator {
	stopCh := make(chan struct{})
	dataCh := make(chan Result, 4)

	calculator := Calculator{
		DataCh: dataCh,

		start:  start,
		step:   big.NewInt(int64(workersCount)),
		stopCh: stopCh,

		cache: NewCache(),
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
			if r, ok := c.cache.Get(number); ok {
				c.DataCh <- r

			} else{
				r = NewResult(number)
				c.DataCh <- r
				c.cache.Put(number, &r)
			}

			number.Add(number, c.step)

		case <-c.stopCh:
			return
		}
	}
}
