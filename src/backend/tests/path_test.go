package tests

import (
	"testing"
	"../calc"
	"math/big"
	"strings"
)

func toBigIntSlice(strSlice []string) ([]*big.Int, bool) {
	var bigIntSlice []*big.Int
	for i := 0; i < len(strSlice); i++ {
		bigInt, success := new(big.Int).SetString(strSlice[i], 10)
		if !success {
			return nil, success
		}
		bigIntSlice = append(bigIntSlice, bigInt)
	}
	return bigIntSlice, true
}

func bigIntSliceToString(a []*big.Int) string {
	var strs []string

	strs = append(strs, "[")
	strs = append(strs, a[0].String())
	for i := 1; i < len(a); i++ {
		strs = append(strs, ", ")
		strs = append(strs, a[i].String())
	}
	strs = append(strs, "]")

	return strings.Join(strs, "")
}

func printExpectedGotBigintSlices(t *testing.T, a, b []*big.Int) {
	t.Error("assertion a == b failed")
	t.Error("expected length ", len(a), " got lenght ", len(b))
	t.Error("expected: ", bigIntSliceToString(a))
	t.Error("got:      ", bigIntSliceToString(b))
}

func assertEqualBigintSlices(t *testing.T, a, b []*big.Int) {
	if len(a) != len(b) {
		printExpectedGotBigintSlices(t, a, b)
		t.Fatal()
	}

	for i := 0; i < len(a); i++ {
		if a[i].Cmp(b[i]) != 0 {
			printExpectedGotBigintSlices(t, a, b)
			t.Fatal()
		}
	}
}

func assertEqual(t *testing.T, a, b *big.Int) {
	if a.Cmp(b) != 0 {
		t.Error("expected: ", a.String())
		t.Error("got:      ", b.String())
		t.Fatal()
	}
}

func TestSimplePath(t *testing.T) {
	request, _ := new(big.Int).SetString("12", 10)
	answerStr := []string {"12", "6", "3", "10", "5", "16", "8", "4", "2", "1"}
	answer, success := toBigIntSlice(answerStr)
	if !success {
		t.Fatal()
	}

	result := calc.FindPath(request)

	assertEqualBigintSlices(t, answer, result)
}

func TestFirstAndLastInPath(t *testing.T) {
	seed := "9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"
	request, _ := new(big.Int).SetString(seed, 10)
	rightFirst := request
	rightLast, _ := new(big.Int).SetString("1", 10)

	result := calc.FindPath(request)

	assertEqual(t, result[0], rightFirst)
	assertEqual(t, result[len(result) - 1], rightLast)
}

func TestOne(t *testing.T) {
	seed := new(big.Int).SetUint64(1)
	result := calc.FindPath(seed)
	one, _ := new(big.Int).SetString("1", 10)

	if len(result) != 1 {
		t.Fatal("Sequence is longer then 1")
	}

	assertEqual(t, result[0], one)
}

func BenchmarkGoLongPath(b* testing.B) {
	// setup test
	seed := "9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"
	request, _ := new(big.Int).SetString(seed, 10)
	rightLength := 3275

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := calc.FindPath(request)
		if len(result) != rightLength {
			b.Fatal()
		}
	}
}

func BenchmarkGoShortPath(b* testing.B) {
	// setup test
	seed := "123456789"
	request, _ := new(big.Int).SetString(seed, 10)
	rightLength := 178

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := calc.FindPath(request)
		if len(result) != rightLength {
			b.Fatal()
		}
	}
}