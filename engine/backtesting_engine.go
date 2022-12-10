package engine

import (
    "fmt"
    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/chart"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/enum"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "github.com/zhengow/vngo/utils"
    "math"
    "sort"
    "time"
)

type BacktestingEngine struct {
    symbols     []*strategy.Symbol
    interval    types.Interval
    start       *strategy.VnTime
    end         *strategy.VnTime
    rates       map[strategy.Symbol]float64
    strategy    strategy.Strategy
    _dts        mapset.Set
    dts         []strategy.VnTime
    historyData map[string]map[strategy.VnTime]strategy.Bar
    datetime    *strategy.VnTime
    tradeCount  int
    *accountEngine
    *statisticEngine
}

var _BacktestingEngine *BacktestingEngine

func NewBacktestingEngine() *BacktestingEngine {
    if _BacktestingEngine != nil {
        return _BacktestingEngine
    }
    _BacktestingEngine = &BacktestingEngine{
        _dts:            mapset.NewSet(),
        historyData:     make(map[string]map[strategy.VnTime]strategy.Bar),
        accountEngine:   newOrderEngine(),
        statisticEngine: newStatisticEngine(),
    }
    return _BacktestingEngine
}

func (b *BacktestingEngine) SetParameters(
    symbols []*strategy.Symbol,
    interval types.Interval,
    start,
    end time.Time,
    rates map[strategy.Symbol]float64,
    priceTicks map[string]int,
    capital float64,
) {
    b.symbols = symbols
    b.interval = interval
    b.start = strategy.NewVnTime(start)
    b.end = strategy.NewVnTime(end)
    b.setRates(rates)
    b.setPriceTicks(priceTicks)
    b.setCapital(capital)
    b.AddCash(capital)
}

func (b *BacktestingEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    strategy.Inject(b.accountEngine)
    b.strategy = strategy
}

func (b *BacktestingEngine) LoadData() {
    defer utils.TimeCost("load data")()
    if b.start == nil || b.end == nil {
        fmt.Println("please set start && end time")
        return
    }
    start := b.start.Format()
    end := b.end.Format()
    for _, symbol := range b.symbols {
        if b.historyData[symbol.Name] == nil {
            b.historyData[symbol.Name] = make(map[strategy.VnTime]strategy.Bar)
        }
        bars := database.LoadBarData(*symbol, b.interval, start, end)
        for _, bar := range bars {
            _time := bar.Datetime
            b._dts.Add(_time)
            b.historyData[symbol.Name][_time] = bar
        }
        fmt.Printf("%s.%s load success, length: %d\n", symbol.Name, symbol.Exchange, len(b.historyData[symbol.Name]))
    }
}

func (b *BacktestingEngine) RunBacktesting() {
    b.dts = make([]strategy.VnTime, b._dts.Cardinality())
    cnt := 0
    b._dts.Each(func(ele interface{}) bool {
        b.dts[cnt] = ele.(strategy.VnTime)
        cnt++
        return false
    })
    sort.Slice(b.dts, func(i, j int) bool {
        return b.dts[i].Time.Before(b.dts[j].Time)
    })

    for _, dt := range b.dts {
        b.newBars(dt)
    }
}

func (b *BacktestingEngine) newBars(dt strategy.VnTime) {
    b.datetime = &dt
    bars := make(map[string]strategy.Bar)
    for _, symbol := range b.symbols {
        bars[symbol.Name] = b.historyData[symbol.Name][dt]
    }
    b.crossLimitOrder(bars)
    b.accountEngine.updateCloses(bars)
    b.strategy.OnBars(bars)
    b.updateClose(bars)
}

func (b *BacktestingEngine) crossLimitOrder(bars map[string]strategy.Bar) {
    for _, order := range b.accountEngine.activeLimitOrders {
        bar := bars[order.Symbol.Name]
        longCrossPrice := bar.LowPrice
        shortCrossPrice := bar.HighPrice
        longBestPrice := bar.OpenPrice
        shortBestPrice := bar.OpenPrice

        longCross := order.Direction == enum.DirectionEnum.LONG && order.Price >= longCrossPrice && longCrossPrice > 0
        shortCross := order.Direction == enum.DirectionEnum.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0

        if !longCross && !shortCross {
            continue
        }

        delete(b.accountEngine.activeLimitOrders, order.OrderId)

        b.tradeCount++

        var tradePrice float64
        if longCross {
            tradePrice = math.Min(order.Price, longBestPrice)
        } else {
            tradePrice = math.Max(order.Price, shortBestPrice)
        }

        tradeData := strategy.NewTradeData(
            order.Symbol,
            order.OrderId,
            b.tradeCount,
            order.Direction,
            tradePrice,
            order.Volume,
            *b.datetime,
        )

        incrementPos := order.Volume
        if order.Direction == enum.DirectionEnum.SHORT {
            incrementPos = -order.Volume
        }
        b.accountEngine.updatePositions(order.Symbol, incrementPos, tradePrice)
        b.strategy.UpdateTrade(*tradeData)
        b.trades[b.tradeCount] = tradeData
    }
}

func (b *BacktestingEngine) ShowPNLChart() {
    chart.ChartPNL(b.dts, b.balances, "")
}

func (b *BacktestingEngine) ShowKLineChart() {
    for _, symbol := range b.symbols {
        bars := make([]strategy.Bar, len(b.dts))
        for idx, dt := range b.dts {
            bars[idx] = b.historyData[symbol.Name][dt]
        }
        trades := make([]*strategy.TradeData, 0)
        for _, trade := range b.trades {
            if trade.Symbol.Name == symbol.Name {
                trades = append(trades, trade)
            }
        }
        chart.ChartKLines(b.dts, bars, trades, symbol.Name)
    }
}
