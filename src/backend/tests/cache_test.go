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

func TestMultipleGet(t *testing.T) {
	cache := calc.NewCache()

	number2, _ := new(big.Int).SetString("1", 10)
	number3, _ := new(big.Int).SetString("10", 10)
	number4, _ := new(big.Int).SetString("25", 10)

	cache.Put(number2, &calc.Result{})
	cache.Put(number3, &calc.Result{})

	number, _ := new(big.Int).SetString("12", 10)
	goodResult := calc.NewResult(number)
	cache.Put(number, &goodResult)

	cache.Put(number4, &calc.Result{})

	realResult, ok := cache.Get(number)

	if !ok{
		t.Fail()
	}

	if !reflect.DeepEqual(goodResult, realResult){
		t.Fail()
	}
}