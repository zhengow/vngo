package main

import "github.com/zhengow/vngo/engine"

func main() {
	b := engine.BacktestingEngine{}
	b.SetParameters()
	b.AddStrategy()
	b.LoadData()
	b.RunBacktesting()
	b.CalculateResult()
}
