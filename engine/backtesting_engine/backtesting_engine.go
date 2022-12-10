package backtesting_engine

import (
    "fmt"
    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/chart"
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "github.com/zhengow/vngo/utils"
    "math"
    "sort"
    "time"
)

type BacktestingEngine struct {
    symbols     []strategy.Symbol
    interval    types.Interval
    start       time.Time
    end         time.Time
    strategy    strategy.Strategy
    _dts        mapset.Set
    dts         []time.Time
    historyData map[string]map[time.Time]strategy.Bar
    datetime    time.Time
    *backtestingAccount
    *statisticEngine
}

func NewBacktestingEngine() *BacktestingEngine {
    return &BacktestingEngine{
        _dts:               mapset.NewSet(),
        historyData:        make(map[string]map[time.Time]strategy.Bar),
        backtestingAccount: newAccount(),
        statisticEngine:    newStatisticEngine(),
    }
}

func (b *BacktestingEngine) AddSymbol(name string, rate float64, exchange types.Exchange, interval types.Interval) *BacktestingEngine {
    symbol := strategy.Symbol{
        Exchange: exchange,
        Name:     name,
        Interval: interval,
    }
    b.setRate(symbol, rate)
    b.symbols = append(b.symbols, symbol)
    return b
}

func (b *BacktestingEngine) AddSymbols(names []string, rates []float64, exchange types.Exchange, interval types.Interval) *BacktestingEngine {
    if len(names) != len(rates) {
        fmt.Println("add failed, len(names) != len(rates)")
        return b
    }
    for idx, name := range names {
        b.AddSymbol(name, rates[idx], exchange, interval)
    }
    return b
}

func (b *BacktestingEngine) StartDate(date time.Time) *BacktestingEngine {
    b.start = date
    return b
}

func (b *BacktestingEngine) EndDate(date time.Time) *BacktestingEngine {
    b.end = date
    return b
}

func (b *BacktestingEngine) Capital(capital float64) *BacktestingEngine {
    b.setCapital(capital)
    b.AddCash(capital)
    return b
}

func (b *BacktestingEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    strategy.Inject(b.backtestingAccount)
    b.strategy = strategy
}

func (b *BacktestingEngine) LoadData() {
    defer utils.TimeCost("load data")()
    if b.start.IsZero() || b.end.IsZero() {
        fmt.Println("please set start && end time")
        return
    }
    start := b.start.Format(consts.DateFormat)
    end := b.end.Format(consts.DateFormat)
    for _, symbol := range b.symbols {
        if b.historyData[symbol.Name] == nil {
            b.historyData[symbol.Name] = make(map[time.Time]strategy.Bar)
        }
        bars := database.LoadBarData(symbol, symbol.Interval, start, end)
        for _, bar := range bars {
            _time := bar.Datetime.Time
            b._dts.Add(_time)
            b.historyData[symbol.Name][_time] = bar
        }
        fmt.Printf("%s.%s load success, length: %d\n", symbol.Name, symbol.Exchange, len(b.historyData[symbol.Name]))
    }
}

func (b *BacktestingEngine) RunBacktesting() {
    b.dts = make([]time.Time, b._dts.Cardinality())
    cnt := 0
    b._dts.Each(func(ele interface{}) bool {
        b.dts[cnt] = ele.(time.Time)
        cnt++
        return false
    })
    sort.Slice(b.dts, func(i, j int) bool {
        return b.dts[i].Before(b.dts[j])
    })

    for _, dt := range b.dts {
        b.newBars(dt)
    }
}

func (b *BacktestingEngine) newBars(dt time.Time) {
    b.datetime = dt
    bars := make(map[string]strategy.Bar)
    for _, symbol := range b.symbols {
        bars[symbol.Name] = b.historyData[symbol.Name][dt]
    }
    b.crossLimitOrder(bars)
    b.backtestingAccount.updateCloses(bars)
    b.strategy.OnBars(bars)
    b.statisticEngine.onBars(bars)
}

func (b *BacktestingEngine) crossLimitOrder(bars map[string]strategy.Bar) {
    for _, order := range b.backtestingAccount.Orders {
        bar := bars[order.Symbol.Name]
        longCrossPrice := bar.LowPrice
        shortCrossPrice := bar.HighPrice
        longBestPrice := bar.OpenPrice
        shortBestPrice := bar.OpenPrice

        if order.Status == consts.StatusEnum.SUBMITTING {
            order.Status = consts.StatusEnum.NOTTRADED
            b.strategy.UpdateOrder(order)
            b.statisticEngine.updateOrder(order)
        }

        longCross := order.Direction == consts.DirectionEnum.LONG && order.Price >= longCrossPrice && longCrossPrice > 0
        shortCross := order.Direction == consts.DirectionEnum.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0

        if !longCross && !shortCross {
            continue
        }

        order.Traded = order.Volume
        order.Status = consts.StatusEnum.ALLTRADED
        b.strategy.UpdateOrder(order)

        delete(b.backtestingAccount.Orders, order.OrderId)

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
            *strategy.NewVnTime(b.datetime),
        )

        incrementPos := order.Volume
        if order.Direction == consts.DirectionEnum.SHORT {
            incrementPos = -order.Volume
        }
        b.backtestingAccount.updatePositions(order.Symbol, incrementPos, tradePrice)
        b.strategy.UpdateTrade(*tradeData)
        b.trades[b.tradeCount] = tradeData
    }
}

func (b *BacktestingEngine) ShowPNLChart() {
    dts := make([]strategy.VnTime, len(b.dts))
    for idx, dt := range b.dts {
        dts[idx] = *strategy.NewVnTime(dt)
    }
    chart.ChartPNL(dts, b.balances, "")
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
        dts := make([]strategy.VnTime, len(b.dts))
        for idx, dt := range b.dts {
            dts[idx] = *strategy.NewVnTime(dt)
        }
        chart.ChartKLines(dts, bars, trades, symbol.Name)
    }
}
