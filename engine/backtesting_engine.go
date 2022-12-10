package engine

import (
    "fmt"
    "github.com/zhengow/vngo"
    "github.com/zhengow/vngo/utils"
    "math"
    "sort"
    "time"

    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/chart"
)

type BacktestingEngine struct {
    symbols     []*vngo.Symbol
    interval    vngo.Interval
    start       *time.Time
    end         *time.Time
    rates       map[vngo.Symbol]float64
    strategy    vngo.Strategy
    _dts        mapset.Set
    dts         []time.Time
    historyData map[string]map[time.Time]vngo.Bar
    datetime    *time.Time
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
        historyData:     make(map[string]map[time.Time]vngo.Bar),
        accountEngine:   newOrderEngine(),
        statisticEngine: newStatisticEngine(),
    }
    return _BacktestingEngine
}

func (b *BacktestingEngine) SetParameters(
    symbols []*vngo.Symbol,
    interval vngo.Interval,
    start,
    end time.Time,
    rates map[vngo.Symbol]float64,
    priceTicks map[string]int,
    capital float64,
) {
    b.symbols = symbols
    b.interval = interval
    b.start = &start
    b.end = &end
    b.setRates(rates)
    b.setPriceTicks(priceTicks)
    b.setCapital(capital)
    b.AddCash(capital)
}

func (b *BacktestingEngine) AddStrategy(strategy vngo.Strategy, setting map[string]interface{}) {
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
    start := b.start.Format(vngo.DateFormat)
    end := b.end.Format(vngo.DateFormat)
    for _, symbol := range b.symbols {
        if b.historyData[symbol.Name] == nil {
            b.historyData[symbol.Name] = make(map[time.Time]vngo.Bar)
        }
        bars := vngo.LoadBarData(*symbol, b.interval, start, end)
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
    b.datetime = &dt
    bars := make(map[string]vngo.Bar)
    for _, symbol := range b.symbols {
        bars[symbol.Name] = b.historyData[symbol.Name][dt]
    }
    b.crossLimitOrder(bars)
    b.accountEngine.updateCloses(bars)
    b.strategy.OnBars(bars)
    b.updateClose(bars)
}

func (b *BacktestingEngine) crossLimitOrder(bars map[string]vngo.Bar) {
    for _, order := range b.accountEngine.activeLimitOrders {
        bar := bars[order.Symbol.Name]
        longCrossPrice := bar.LowPrice
        shortCrossPrice := bar.HighPrice
        longBestPrice := bar.OpenPrice
        shortBestPrice := bar.OpenPrice

        longCross := order.Direction == vngo.DirectionEnum.LONG && order.Price >= longCrossPrice && longCrossPrice > 0
        shortCross := order.Direction == vngo.DirectionEnum.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0

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

        tradeData := vngo.NewTradeData(
            order.Symbol,
            order.OrderId,
            b.tradeCount,
            order.Direction,
            tradePrice,
            order.Volume,
            *b.datetime,
        )

        incrementPos := order.Volume
        if order.Direction == vngo.DirectionEnum.SHORT {
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
        bars := make([]vngo.Bar, len(b.dts))
        for idx, dt := range b.dts {
            bars[idx] = b.historyData[symbol.Name][dt]
        }
        trades := make([]*vngo.TradeData, 0)
        for _, trade := range b.trades {
            if trade.Symbol.Name == symbol.Name {
                trades = append(trades, trade)
            }
        }
        chart.ChartKLines(b.dts, bars, trades, symbol.Name)
    }
}
