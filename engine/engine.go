package engine

import (
    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "time"
)

type BaseEngine struct {
    symbols     []models.Symbol
    interval    types.Interval
    start       time.Time
    end         time.Time
    strategy    strategy.Strategy
    _dts        mapset.Set
    dts         []time.Time
    historyData map[string]map[time.Time]models.Bar
    datetime    time.Time
    *BaseAccount
}

func (b *BaseEngine) AddSymbol(name string, exchange types.Exchange) *BaseEngine {
    symbol := models.Symbol{
        Exchange: exchange,
        Name:     name,
    }
    b.symbols = append(b.symbols, symbol)
    return b
}

func (b *BaseEngine) AddSymbols(names []string, exchange types.Exchange) *BaseEngine {
    for _, name := range names {
        b.AddSymbol(name, exchange)
    }
    return b
}

func (b *BaseEngine) SetInterval(interval types.Interval) *BaseEngine {
    b.interval = interval
    return b
}

func (b *BaseEngine) StartDate(date time.Time) *BaseEngine {
    b.start = date
    return b
}

func (b *BaseEngine) EndDate(date time.Time) *BaseEngine {
    b.end = date
    return b
}

func (b *BaseEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    strategy.Inject(b.BaseAccount)
    b.strategy = strategy
}

func (b *BaseEngine) LoadData() {
}
