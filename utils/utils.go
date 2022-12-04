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

func RoundTo(price float64, priceTick int) float64 {
    res, _ := strconv.ParseFloat(strconv.FormatFloat(price, 'f', priceTick, 64), 64)
    return res
}
