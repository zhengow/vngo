package main

import (
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/database"
	"github.com/zhengow/vngo/engine"
	"github.com/zhengow/vngo/strategy"
)

func main() {
	b := engine.NewBacktestingEngine()
	b.SetParameters([]string{"BTCDOMUSDT.BINANCE"}, consts.IntervalEnum.MINUTE, time.Now().AddDate(0, 0, -7), time.Now(), nil, 10000)
	b.AddStrategy(strategy.MyStrategy{}, map[string]interface{}{"Test": 1})
	database.NewMysql()
	b.LoadData()
	b.RunBacktesting()
	b.CalculateResult()
}
