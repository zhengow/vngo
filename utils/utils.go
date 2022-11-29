package utils

import (
	"fmt"
	"time"
)

func TimeCost(name string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("%s time cost = %v\n", name, tc)
	}
}
