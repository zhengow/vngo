package main

import (
	"fmt"
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/database"
	"github.com/zhengow/vngo/engine"
	"github.com/zhengow/vngo/model"
	"github.com/zhengow/vngo/strategy"
	"github.com/zhengow/vngo/utils"
)

func main() {
	b := engine.BacktestingEngine{}
	// b.SetParameters()
	b.AddStrategy(strategy.Strategy{}, map[string]interface{}{"Test": 1})
	sqlite := database.NewSqlite()
	bars := make([]model.Bar, 0)
	for i := 0; i < 100; i++ {
		bar := model.Bar{
			Symbol:   "testsymbol",
			Exchange: consts.ExchangeEnum.BINANCE,
			Datetime: utils.DatabaseTime(time.Now()),
			Interval: consts.IntervalEnum.MINUTE,
			Volume:   float64(i),
		}
		bars = append(bars, bar)
	}
	sqlite.SaveBarData(bars)
	b.LoadData(sqlite)
	b.RunBacktesting()
	b.CalculateResult()
}
