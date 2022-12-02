package engine

import (
    "fmt"
    "math"
    "sort"
    "time"

    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/model"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/utils"
)

type BacktestingEngine struct {
    symbols     []string
    interval    consts.Interval
    start       *time.Time
    end         *time.Time
    rates       map[string]float64
    capital     float64
    strategy    strategy.Strategy
    dts         mapset.Set
    historyData map[string]map[time.Time]model.Bar
    datetime    *time.Time
    priceTicks  map[string]int
    tradeCount  int
    trades      map[int]*model.TradeData
    *orderEngine
}

var _BacktestingEngine *BacktestingEngine

func NewBacktestingEngine() *BacktestingEngine {
    if _BacktestingEngine != nil {
        return _BacktestingEngine
    }
    _BacktestingEngine = &BacktestingEngine{
        dts:         mapset.NewSet(),
        historyData: make(map[string]map[time.Time]model.Bar),
        trades:      make(map[int]*model.TradeData),
    }
    return _BacktestingEngine
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
    b.start = &start
    b.end = &end
    b.rates = rates
    b.capital = capital
}

func (b *BacktestingEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(setting)
    b.strategy = strategy
}

func (b *BacktestingEngine) LoadData() {
    defer utils.TimeCost("load data")()
    if b.start == nil || b.end == nil {
        fmt.Println("please set start && end time")
        return
    }
    start := b.start.Format(consts.DateFormat)
    end := b.end.Format(consts.DateFormat)
    for _, _symbol := range b.symbols {
        symbol, exchange := utils.ParseSymbol(_symbol)
        if symbol == "" || exchange == "" {
            continue
        }
        if b.historyData[_symbol] == nil {
            b.historyData[_symbol] = make(map[time.Time]model.Bar)
        }
        bars := database.LoadBarData(symbol, consts.Exchange(exchange), b.interval, start, end)
        for _, bar := range bars {
            _time := time.Time(bar.Datetime)
            b.dts.Add(_time)
            b.historyData[_symbol][_time] = bar
        }
        fmt.Printf("%s load success, length: %d\n", symbol, len(b.historyData[_symbol]))
    }
}

func (b *BacktestingEngine) RunBacktesting() {
    b.orderEngine = newOrderEngine(b.priceTicks)
    b.strategy.Inject(b.orderEngine)
    b.strategy.OnInit()
    dts := make([]time.Time, b.dts.Cardinality())
    cnt := 0
    b.dts.Each(func(ele interface{}) bool {
        dts[cnt] = ele.(time.Time)
        cnt++
        return false
    })
    sort.Slice(dts, func(i, j int) bool {
        return dts[i].Before(dts[j])
    })

    for _, dt := range dts {
        if b.datetime != nil && dt.Day() != b.datetime.Day() {
            b.strategy.DoneInit()
        }
        b.newBars(dt)
    }
}

func (b *BacktestingEngine) newBars(dt time.Time) {
    b.datetime = &dt
    bars := make(map[string]model.Bar)
    for _, symbol := range b.symbols {
        bars[symbol] = b.historyData[symbol][dt]
    }
    b.crossLimitOrder(bars)
    b.strategy.OnBars(bars)
}

func (b *BacktestingEngine) crossLimitOrder(bars map[string]model.Bar) {
    for _, order := range b.orderEngine.activeLimitOrders {
        bar := bars[order.Symbol]
        longCrossPrice := bar.LowPrice
        shortCrossPrice := bar.HighPrice
        longBestPrice := bar.OpenPrice
        shortBestPrice := bar.OpenPrice

        longCross := order.Direction == consts.DirectionEnum.LONG && order.Price >= longCrossPrice && longCrossPrice > 0
        shortCross := order.Direction == consts.DirectionEnum.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0

        if !longCross && !shortCross {
            continue
        }

        delete(b.orderEngine.activeLimitOrders, order.OrderId)

        b.tradeCount++

        var tradePrice float64
        if longCross {
            tradePrice = math.Min(order.Price, longBestPrice)
        } else {
            tradePrice = math.Max(order.Price, shortBestPrice)
        }

        tradeData := model.NewTradeData(
            order.Symbol,
            order.Exchange,
            order.OrderId,
            b.tradeCount,
            order.Direction,
            tradePrice,
            order.Volume,
            *b.datetime,
        )

        b.strategy.UpdateTrade(*tradeData)
        b.trades[b.tradeCount] = tradeData
    }
}

func (b *BacktestingEngine) CalculateResult() {
    println("calc")
}
