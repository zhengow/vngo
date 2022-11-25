package engine

import (
	"time"

	"github.com/zhengow/vngo/strategy"
)

type Interval string

type BacktestingEngine struct {
	symbols  []string
	interval Interval
	start    time.Time
	end      time.Time
	rates    map[string]float64
	capital  float64
	strategy strategy.Strategy
}

func (b *BacktestingEngine) SetParameters(
	symbols []string,
	interval Interval,
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
	println("load")
}

func (b *BacktestingEngine) RunBacktesting() {
	println("run")
}

func (b *BacktestingEngine) CalculateResult() {
	println("calc")
}
