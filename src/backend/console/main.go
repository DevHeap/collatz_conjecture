package main

import (
	"fmt"
	"math/big"
	"time"

	"../calc"
)

// Simple console interface for testing backend logic

func main() {
	fmt.Println("Enter starting integer")
	var input string
	fmt.Scanln(&input)

	number, parsed := new(big.Int).SetString(input, 10)
	if !parsed {
		panic("Unable to parse string into integer number with base 10")
	}

	calculator := calc.NewCalculator(number, 4)

	for {
		r := <-calculator.DataCh
		fmt.Println(
			"\tNumber: ", r.Number,
			"\tPath Len:", r.PathLength,
			"\tElapsed Time: ", r.Time)

		time.Sleep(time.Millisecond * 500)
	}
}
