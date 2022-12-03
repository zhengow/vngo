package utils

import (
    "fmt"
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

//
//func ParseSymbol(symbol string) (string, string) {
//	parts := strings.Split(symbol, ".")
//	if len(parts) == 2 {
//		return parts[0], parts[1]
//	}
//	fmt.Println("invalid symbol: ", symbol)
//	return "", ""
//}

func RoundTo(price float64, priceTick int) float64 {
    res, _ := strconv.ParseFloat(strconv.FormatFloat(price, 'f', priceTick, 64), 64)
    return res
}

func FindSmaller(price float64, prices []float64) float64 {
    reversedPrices := make([]float64, len(prices))
    for i := 0; i < len(prices); i++ {
        reversedPrices[i] = prices[len(prices)-i-1]
    }
    for _, p := range reversedPrices {
        if p < price {
            return p
        }
    }
    return 0
}

func FindLarger(price float64, prices []float64) float64 {
    for _, p := range prices {
        if p > price {
            return p
        }
    }
    return prices[len(prices)-1]
}
