package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeCost(name string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("%s time cost = %v\n", name, tc)
	}
}

func ParseSymbol(symbol string) (string, string) {
	parts := strings.Split(symbol, ".")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	fmt.Println("invalid symbol: ", symbol)
	return "", ""
}

func RoundTo(price float64, priceTick int) float64 {
	res, _ := strconv.ParseFloat(strconv.FormatFloat(price, 'f', priceTick, 64), 64)
	return res
}
