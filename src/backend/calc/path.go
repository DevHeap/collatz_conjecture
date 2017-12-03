package calc

import (
	"math/big"
	"errors"
)

var zero  = big.NewInt(0)
var one   = big.NewInt(1)
var two   = big.NewInt(2)
var three = big.NewInt(3)


func step(x *big.Int) *big.Int{
	result := new(big.Int)

	mod := new(big.Int).Mod(x, two)

	if mod.Cmp(zero) == 0 {
		result.Div(x, two)
	}else{
		result.Mul(x,    three)
		result.Add(result, one)
	}

	return result
}

func FindPath(x *big.Int) []*big.Int {
	var path []*big.Int


	for x.Cmp(one) > 0{
		x = step(x)
		path = append(path, x)
	}

	return path
}

type result struct{
	number *big.Int
	path   []*big.Int
}

type worker struct{
	number *big.Int
	step   *big.Int

	stopCh chan struct{}
	dataCh chan result
}


func (w *worker) compute(){
	for {
		select{
		default:
			w.dataCh <- result {
				number: w.number,
				path: FindPath(w.number),
			}

			w.number.Add(w.number, w.step)

		case <-w.stopCh:
			return
		}
	}
}

func startWorkers(start *big.Int, count int){
	stopCh := make(chan  struct{})
	dataCh := make(chan  result)

	

}

//number, parsed := new(big.Int).SetString(str_number, 10)
//if !parsed {
//	return nil, errors.New("Unable to parse string into integer number with base 10")
//}

func FindPathsStartingFrom(start *big.Int, step *big.Int, stop_ch chan struct{}){
	number := new(big.Int).Set(start)
	for {
		select{
		default:
			path := FindPath(number)
			number.Add(number, step)

		case <-stop_ch:
			return
		}
	}
}

