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

func ParseBarData(open, high, low, close, volume string) (float64, float64, float64, float64, float64, error) {
	var err error
	var openValue, highValue, lowValue, closeValue, volumeValue float64
	openValue, err = strconv.ParseFloat(open, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	highValue, err = strconv.ParseFloat(high, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	lowValue, err = strconv.ParseFloat(low, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	closeValue, err = strconv.ParseFloat(close, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	volumeValue, err = strconv.ParseFloat(volume, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	return openValue, highValue, lowValue, closeValue, volumeValue, nil
}
