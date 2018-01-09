package tests

import (
	"testing"
	"math/big"
	"../calc"
	"reflect"
)

func TestSimpleGet(t *testing.T) {
	cache := calc.NewCache()

	number, _ := new(big.Int).SetString("12", 10)
	goodResult := calc.NewResult(number)

	cache.Put(number, &goodResult)
	realResult, ok := cache.Get(number)

	if !ok{
		t.Fail()
	}

	if !reflect.DeepEqual(goodResult, realResult){
		t.Fail()
	}
}