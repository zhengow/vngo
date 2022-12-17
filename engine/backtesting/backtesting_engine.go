package backtesting

import (
    "fmt"
    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/chart"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/utils"
    "math"
    "sort"
    "time"
)

type Engine struct {
    _dts        mapset.Set
    dts         []time.Time
    historyData map[string]map[time.Time]models.Bar
    datetime    time.Time
    *backtestingAccount
    *statisticEngine
    engine.BaseEngine
}

func NewBacktestingEngine() *Engine {
    return &Engine{
        _dts:               mapset.NewSet(),
        historyData:        make(map[string]map[time.Time]models.Bar),
        backtestingAccount: newAccount(),
        statisticEngine:    newStatisticEngine(),
    }
}

func (b *Engine) Capital(capital float64) *Engine {
    b.setCapital(capital)
    b.AddCash(capital)
    return b
}

func (b *Engine) LoadHistoryData(start, end time.Time) {
    defer utils.TimeCost("load history data")()
    if start.IsZero() || end.IsZero() {
        fmt.Println("please set start && end time")
        return
    }
    for _, symbol := range b.GetSymbols() {
        if b.historyData[symbol.FullName()] == nil {
            b.historyData[symbol.FullName()] = make(map[time.Time]models.Bar)
        }
        bars := database.LoadBarData(symbol, b.GetInterval(), start.Format(models.DateFormat), end.Format(models.DateFormat))
        for _, bar := range bars {
            _time := bar.Datetime.Time
            b._dts.Add(_time)
            b.historyData[symbol.FullName()][_time] = bar
        }
        fmt.Printf("%s.%s load success, length: %d\n", symbol.FullName(), symbol.Exchange, len(b.historyData[symbol.FullName()]))
    }
}

func (b *Engine) RunBacktesting() {
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
        //fmt.Printf("%d/%d\n", idx, len(b.dts))
        b.newBars(dt)
    }
}

func (b *Engine) newBars(dt time.Time) {
    b.datetime = dt
    bars := make(map[string]models.Bar)
    for _, symbol := range b.GetSymbols() {
        bars[symbol.FullName()] = b.historyData[symbol.FullName()][dt]
    }
    b.crossLimitOrder(bars)
    b.backtestingAccount.updateCloses(bars)
    b.Queue.Bars <- bars
    <-b.Queue.Continue
    b.statisticEngine.onBars(bars)
}

func (b *Engine) crossLimitOrder(bars map[string]models.Bar) {
    for _, order := range b.backtestingAccount.Orders {
        bar := bars[order.Symbol.FullName()]
        longCrossPrice := bar.LowPrice
        shortCrossPrice := bar.HighPrice
        longBestPrice := bar.OpenPrice
        shortBestPrice := bar.OpenPrice

        if order.Status == models.StatusEnum.SUBMITTING {
            order.Status = models.StatusEnum.NOTTRADED
            b.Queue.Order <- order
            <-b.Queue.Continue
            b.statisticEngine.updateOrder(order)
        }

        longCross := order.Direction == models.DirectionEnum.LONG && order.Price >= longCrossPrice && longCrossPrice > 0
        shortCross := order.Direction == models.DirectionEnum.SHORT && order.Price <= shortCrossPrice && shortCrossPrice > 0

        if !longCross && !shortCross {
            continue
        }

        order.Traded = order.Volume
        order.Status = models.StatusEnum.ALLTRADED

        b.Queue.Order <- order
        <-b.Queue.Continue

        delete(b.backtestingAccount.Orders, order.OrderId)

        b.tradeCount++

        var tradePrice float64
        if longCross {
            tradePrice = math.Min(order.Price, longBestPrice)
        } else {
            tradePrice = math.Max(order.Price, shortBestPrice)
        }

        tradeData := models.NewTradeData(
            order.Symbol,
            order.OrderId,
            b.tradeCount,
            order.Direction,
            tradePrice,
            order.Volume,
            models.NewVnTime(b.datetime),
        )

        incrementPos := order.Volume
        if order.Direction == models.DirectionEnum.SHORT {
            incrementPos = -order.Volume
        }
        b.backtestingAccount.updatePositions(order.Symbol, incrementPos, tradePrice)
        b.Queue.Trade <- *tradeData
        <-b.Queue.Continue
        b.trades[b.tradeCount] = tradeData
    }
}

func (b *Engine) ShowPNLChart() {
    dts := make([]models.VnTime, len(b.dts))
    for idx, dt := range b.dts {
        dts[idx] = models.NewVnTime(dt)
    }
    chart.ChartPNL(dts, b.balances, "")
}

func (b *Engine) ShowKLineChart() {
    for _, symbol := range b.GetSymbols() {
        bars := make([]models.Bar, len(b.dts))
        for idx, dt := range b.dts {
            bars[idx] = b.historyData[symbol.FullName()][dt]
        }
        trades := make([]*models.TradeData, 0)
        for _, trade := range b.trades {
            if trade.Symbol.FullName() == symbol.FullName() {
                trades = append(trades, trade)
            }
        }
        dts := make([]models.VnTime, len(b.dts))
        for idx, dt := range b.dts {
            dts[idx] = models.NewVnTime(dt)
        }
        chart.ChartKLines(dts, bars, trades, symbol.FullName())
    }
}

func (b *Engine) GetAccount() models.Account {
    return b.backtestingAccount
}
