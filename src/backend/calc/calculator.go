package calc

import (
	"math/big"
	"time"
)

type Result struct {
	Number *big.Int
	Path   []*big.Int
	Time   time.Duration
	Offset int64
}

type Calculator struct {
	DataCh chan Result

	start  *big.Int
	step   *big.Int
	stopCh chan struct{}
}

func NewCalculator(start *big.Int, workersCount int) Calculator {
	stopCh := make(chan struct{})
	dataCh := make(chan Result, 4)

	calculator := Calculator{
		DataCh: dataCh,

		start:  start,
		step:   big.NewInt(int64(workersCount)),
		stopCh: stopCh,
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
			start := time.Now()
			path := FindPath(number)
			elapsed := time.Since(start)

			c.DataCh <- Result{
				Number: new(big.Int).Set(number),
				Path:   path,
				Time:   elapsed,
				Offset: offset,
			}

			number.Add(number, c.step)

		case <-c.stopCh:
			return
		}
	}
}
