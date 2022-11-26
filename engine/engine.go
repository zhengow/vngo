package engine

import (
	"fmt"
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/database"
	"github.com/zhengow/vngo/strategy"
)

type BacktestingEngine struct {
	symbols  []string
	interval consts.Interval
	start    time.Time
	end      time.Time
	rates    map[string]float64
	capital  float64
	strategy strategy.Strategy
}

func (b *BacktestingEngine) SetParameters(
	symbols []string,
	interval consts.Interval,
	start,
	end time.Time,
	rates map[string]float64,
	capital float64,
) {
	b.symbols = symbols
	b.interval = interval
	b.start = start
	b.end = end
	b.rates = rates
	b.capital = capital
}

func (b *BacktestingEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
	strategy.SetSetting(setting)
	b.strategy = strategy
}

func (b *BacktestingEngine) LoadData() {
	start := time.Now().AddDate(-1, 0, 0).Format(consts.DateFormat)
	end := time.Now().Format(consts.DateFormat)
	bars := database.LoadBarData("BTCDOMUSDT", consts.ExchangeEnum.BINANCE, consts.IntervalEnum.MINUTE, start, end)
	for _, bar := range bars {
		fmt.Println(time.Time(bar.Datetime).Format(consts.DateFormat))
	}
}

func (b *BacktestingEngine) RunBacktesting() {
	println("run")
}

func (b *BacktestingEngine) CalculateResult() {
	println("calc")
}
