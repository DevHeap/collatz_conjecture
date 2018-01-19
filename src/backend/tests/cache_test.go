package tests

import (
	"testing"
	"math/big"
	"../calc/cache"
	"../calc"
	"reflect"
)

func TestSimpleGet(t *testing.T) {
	// Simple test for base functionality of Get and Put

	cache := cache.NewCacheDefault()

	number, _ := new(big.Int).SetString("42", 10)
	storedPath := calc.FindPath(number)

	cache.Put(storedPath)
	realResult, ok := cache.Get(number)

	if !ok{
		t.Fail()
	}
	if !reflect.DeepEqual(storedPath, realResult){
		t.Fail()
	}
}

func TestMultipleGet(t *testing.T) {
	// Now we adding path and try to get number in this path (for sub path)
	// For 42 path is 42 > 21 > 64 > 32 > 16 > 8 > 4 > 2 > 1

	cache := cache.NewCacheDefault()

	number, _ := new(big.Int).SetString("42", 10)

	fullPath := calc.FindPath(number)
	cache.Put(fullPath)

	subPath  := []*big.Int{
		big.NewInt(8),
		big.NewInt(4),
		big.NewInt(2),
		big.NewInt(1)}

	result, ok := cache.Get(big.NewInt(8))

	if !ok{
		t.Fail()
	}

	if !reflect.DeepEqual(result, subPath){
		t.Fail()
	}
}