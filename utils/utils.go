package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func TimeCost(name string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("%s time cost = %v\n", name, tc)
	}
}

func RoundTo(price float64, priceTick int) float64 {
	res, _ := strconv.ParseFloat(strconv.FormatFloat(price, 'f', priceTick, 64), 64)
	return res
}

func YearMinutes() float64 {
	return (time.Hour * 24 * 365).Minutes()
}

func Mean(v []float64) float64 {
	var res float64 = 0
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += v[i]
	}
	return res / float64(n)
}

func Variance(v []float64) float64 {
	var res float64 = 0
	var m = Mean(v)
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += (v[i] - m) * (v[i] - m)
	}
	return res / float64(n-1)
}
func Std(v []float64) float64 {
	return math.Sqrt(Variance(v))
}
