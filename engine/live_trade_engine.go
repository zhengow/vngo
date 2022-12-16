package engine

import (
    "fmt"
    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/gateway"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "github.com/zhengow/vngo/utils"
    "log"
    "sort"
)

type LiveTradeEngine struct {
    BaseEngine
    symbols     []models.Symbol
    interval    types.Interval
    strategy    strategy.Strategy
    datetime    *models.VnTime
    _dts        mapset.Set
    historyData map[string]map[models.VnTime]models.Bar
    account     models.Account
    gI          gateway.GatewayInterface
}

var _LiveTradeEngine *LiveTradeEngine

func NewLiveTradeEngine(gI gateway.GatewayInterface) *LiveTradeEngine {
    if _LiveTradeEngine != nil {
        return _LiveTradeEngine
    }
    _LiveTradeEngine = &LiveTradeEngine{
        _dts:        mapset.NewSet(),
        account:     nil,
        gI:          gI,
        historyData: make(map[string]map[models.VnTime]models.Bar),
    }
    return _LiveTradeEngine
}

func (b *LiveTradeEngine) SetParameters(
    symbols []models.Symbol,
    interval types.Interval,
) {
    b.symbols = symbols
    b.interval = interval
}

func (b *LiveTradeEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    //models.Inject(b.AccountEngine)
    b.strategy = strategy
}

func (b *LiveTradeEngine) LoadData() {
    defer utils.TimeCost("load data")()
    for _, symbol := range b.symbols {
        bars, err := b.gI.LoadBarData(symbol, b.interval)
        if err != nil {
            log.Fatal(err)
            panic(err)
        }
        if bars != nil {
            if b.historyData[symbol.Name] == nil {
                b.historyData[symbol.Name] = make(map[models.VnTime]models.Bar)
            }
            for _, bar := range bars {
                _time := bar.Datetime
                b._dts.Add(_time)
                b.historyData[symbol.Name][_time] = bar
            }
            fmt.Printf("%s.%s load success, length: %d\n", symbol.Name, symbol.Exchange, len(b.historyData[symbol.Name]))
        }
    }
    b.preRun()
}

func (b *LiveTradeEngine) preRun() {
    dts := make([]models.VnTime, b._dts.Cardinality())
    cnt := 0
    b._dts.Each(func(ele interface{}) bool {
        dts[cnt] = ele.(models.VnTime)
        cnt++
        return false
    })
    sort.Slice(dts, func(i, j int) bool {
        return dts[i].Time.Before(dts[j].Time)
    })

    for _, dt := range dts {
        b.newBars(dt)
    }
}

func (b *LiveTradeEngine) Run() {
    b.gI.WebSocketKLine(b.symbols, b.interval)
    // b.vngo.
    // dts = make([]models.VnTime, b._dts.Cardinality())
    // cnt := 0
    // b._dts.Each(func(ele interface{}) bool {
    // 	dts[cnt] = ele.(models.VnTime)
    // 	cnt++
    // 	return false
    // })
    // sort.Slice(dts, func(i, j int) bool {
    // 	return dts[i].Before(dts[j])
    // })

    // for _, dt := range dts {
    // 	b.newBars(dt)
    // }
}

func (b *LiveTradeEngine) newBars(dt models.VnTime) {
    bars := make(map[models.Symbol]models.Bar)
    for _, symbol := range b.symbols {
        bars[symbol] = b.historyData[symbol.Name][dt]
    }
    b.strategy.OnBars(bars)
}
