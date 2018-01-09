package calc

import (
	"math/big"
	"time"
)

type Result struct {
	Number        *big.Int      // n (seed integer)
	PathLength    int           // path length
	MaxNumber     *big.Int      // highest number
	AverageNumber *big.Int      // average number in path
	Time          time.Duration // calculation time
}

func NewResult(seed *big.Int) Result {
	start := time.Now()
	path := FindPath(seed)
	elapsed := time.Since(start)

	// In case of zero path length (numbers from -inf to 1)
	// Just return bare minimum info
	if len(path) == 0 {
		return Result{
			Number: new(big.Int).Set(seed),
			PathLength: 0,
			Time: elapsed,
		}
	}

	max    := path[0] //path cannot be empty
	length := len(path)
	total  := big.NewInt(0)

	for _, number := range path {
		//if number > max
		if number.Cmp(max) >= 1 {
			max = number
		}

		total.Add(total, number)
	}

	average := big.NewInt(0).Div(total, big.NewInt(int64(length)))

	return Result{
		Number:        new(big.Int).Set(seed),
		PathLength:    length,
		MaxNumber:     max,
		AverageNumber: average,
		Time:          elapsed,
	}
}