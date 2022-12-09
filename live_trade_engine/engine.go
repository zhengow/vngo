package live_trade_engine

import (
    "fmt"
    "math"
    "sort"
    "time"

    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/chart"
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/model"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/utils"
)

type LiveTradeEngine struct {
    symbols     []*model.Symbol
    interval    consts.Interval
    start       *time.Time
    end         *time.Time
    rates       map[model.Symbol]float64
    strategy    strategy.Strategy
    _dts        mapset.Set
    dts         []time.Time
    historyData map[string]map[time.Time]model.Bar
    datetime    *time.Time
    tradeCount  int
    *accountEngine
    *statisticEngine
}

var _LiveTradeEngine *LiveTradeEngine

func NewEngine() *LiveTradeEngine {
    if _LiveTradeEngine != nil {
        return _LiveTradeEngine
    }
    _LiveTradeEngine = &LiveTradeEngine{
        _dts:            mapset.NewSet(),
        historyData:     make(map[string]map[time.Time]model.Bar),
        accountEngine:   newOrderEngine(),
        statisticEngine: newStatisticEngine(),
    }
    return _LiveTradeEngine
}

func (b *LiveTradeEngine) SetParameters(
    symbols []*model.Symbol,
    interval consts.Interval,
    start,
    end time.Time,
    rates map[model.Symbol]float64,
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

func (b *LiveTradeEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    strategy.Inject(b.accountEngine)
    b.strategy = strategy
}

func (b *LiveTradeEngine) LoadData() {
    defer utils.TimeCost("load data")()
    if b.start == nil || b.end == nil {
        fmt.Println("please set start && end time")
        return
    }
    start := b.start.Format(consts.DateFormat)
    end := b.end.Format(consts.DateFormat)
    for _, symbol := range b.symbols {
        if b.historyData[symbol.Symbol] == nil {
            b.historyData[symbol.Symbol] = make(map[time.Time]model.Bar)
        }
        bars := database.LoadBarData(*symbol, b.interval, start, end)
        for _, bar := range bars {
            _time := time.Time(bar.Datetime)
            b._dts.Add(_time)
            b.historyData[symbol.Symbol][_time] = bar
        }
        fmt.Printf("%s.%s load success, length: %d\n", symbol.Symbol, symbol.Exchange, len(b.historyData[symbol.Symbol]))
    }
}

func (b *LiveTradeEngine) RunBacktesting() {
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

func (b *LiveTradeEngine) newBars(dt time.Time) {
    b.datetime = &dt
    bars := make(map[string]model.Bar)
    for _, symbol := range b.symbols {
        bars[symbol.Symbol] = b.historyData[symbol.Symbol][dt]
    }
    b.crossLimitOrder(bars)
    b.accountEngine.updateCloses(bars)
    b.strategy.OnBars(bars)
    b.updateClose(bars)
}

func (b *LiveTradeEngine) crossLimitOrder(bars map[string]model.Bar) {
    for _, order := range b.accountEngine.activeLimitOrders {
        bar := bars[order.Symbol.Symbol]
        longCrossPrice := bar.LowPrice
        shortCrossPrice := bar.HighPrice
        longBestPrice := bar.OpenPrice
        shortBestPrice := bar.OpenPrice

        longCross := order.Direction == consts.DirectionEnum.LONG && order.Price >= longCrossPrice && longCrossPrice > 0
        shortCross := order.Direction == consts.DirectionEnum.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0

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

        tradeData := model.NewTradeData(
            order.Symbol,
            order.OrderId,
            b.tradeCount,
            order.Direction,
            tradePrice,
            order.Volume,
            *b.datetime,
        )

        incrementPos := order.Volume
        if order.Direction == consts.DirectionEnum.SHORT {
            incrementPos = -order.Volume
        }
        b.accountEngine.updatePositions(order.Symbol, incrementPos, tradePrice)
        b.strategy.UpdateTrade(*tradeData)
        b.trades[b.tradeCount] = tradeData
    }
}

func (b *LiveTradeEngine) ShowPNLChart() {
    chart.ChartPNL(b.dts, b.balances, "")
}

func (b *LiveTradeEngine) ShowKLineChart() {
    for _, symbol := range b.symbols {
        bars := make([]model.Bar, len(b.dts))
        for idx, dt := range b.dts {
            bars[idx] = b.historyData[symbol.Symbol][dt]
        }
        trades := make([]*model.TradeData, 0)
        for _, trade := range b.trades {
            if trade.Symbol.Symbol == symbol.Symbol {
                trades = append(trades, trade)
            }
        }
        chart.ChartKLines(b.dts, bars, trades, symbol.Symbol)
    }
}