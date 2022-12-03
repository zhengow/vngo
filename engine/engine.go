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
    symbols     []*model.Symbol
    interval    consts.Interval
    start       *time.Time
    end         *time.Time
    rates       map[string]float64
    capital     float64
    strategy    strategy.Strategy
    _dts        mapset.Set
    dts         []time.Time
    historyData map[string]map[time.Time]model.Bar
    datetime    *time.Time
    priceTicks  map[string]int
    tradeCount  int
    trades      map[int]*model.TradeData
    closes      map[time.Time]map[string]float64
    *orderEngine
}

var _BacktestingEngine *BacktestingEngine

func NewBacktestingEngine() *BacktestingEngine {
    if _BacktestingEngine != nil {
        return _BacktestingEngine
    }
    _BacktestingEngine = &BacktestingEngine{
        _dts:        mapset.NewSet(),
        historyData: make(map[string]map[time.Time]model.Bar),
        trades:      make(map[int]*model.TradeData),
        closes:      make(map[time.Time]map[string]float64),
    }
    return _BacktestingEngine
}

func (b *BacktestingEngine) SetParameters(
    symbols []*model.Symbol,
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

func (b *BacktestingEngine) RunBacktesting() {
    b.orderEngine = newOrderEngine(b.priceTicks)
    b.strategy.Inject(b.orderEngine)
    b.strategy.OnInit()
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
        bars[symbol.Symbol] = b.historyData[symbol.Symbol][dt]
    }
    b.crossLimitOrder(bars)
    b.strategy.OnBars(bars)
    b.updateClose(bars)
}

func (b *BacktestingEngine) updateClose(bars map[string]model.Bar) {
    currentCloses := make(map[string]float64)
    for symbol, bar := range bars {
        currentCloses[symbol] = bar.ClosePrice
    }
    b.closes[*b.datetime] = currentCloses
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
    trades := make(map[time.Time][]*model.TradeData)
    for _, trade := range b.trades {
        if dtTrades, ok := trades[trade.Datetime]; ok {
            dtTrades = append(dtTrades, trade)
        } else {
            trades[trade.Datetime] = []*model.TradeData{trade}
        }
    }
    currentPos := make(map[string]float64)
    netPnls := make([]float64, len(b.dts))
    netPnls[0] = b.capital
    for idx, dt := range b.dts[1:] {
        preCloses := b.closes[b.dts[idx]]
        closes := b.closes[dt]
        pnl := netPnls[idx]
        for symbol, _close := range closes {
            pos := currentPos[symbol]
            if pos == 0 {
                continue
            }
            pnl += pos * (_close - preCloses[symbol])
        }
        if dtTrades, ok := trades[dt]; ok {
            for _, _trade := range dtTrades {
                symbol := _trade.Symbol
                volume := _trade.Volume
                if _trade.Direction == consts.DirectionEnum.SHORT {
                    volume *= -1
                }
                currentPos[symbol] += volume
                pnl += volume * (closes[symbol] - _trade.Price)
            }
        }
        netPnls[idx+1] = pnl
    }
    //fmt.Println(netPnls)
    chart(b.dts, netPnls, "")
}
