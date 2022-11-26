package main

import (
	"github.com/zhengow/vngo/database"
	"github.com/zhengow/vngo/engine"
	"github.com/zhengow/vngo/strategy"
)

func main() {
	b := engine.BacktestingEngine{}
	// b.SetParameters()
	b.AddStrategy(strategy.Strategy{}, map[string]interface{}{"Test": 1})
	database.NewMysql()
	b.LoadData()
	b.RunBacktesting()
	b.CalculateResult()
}
