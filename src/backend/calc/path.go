package calc

import (
	"math/big"
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
