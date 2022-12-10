package engine

import (
    "fmt"
    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/gateway"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "github.com/zhengow/vngo/utils"
    "log"
    "sort"
)

type LiveTradeEngine struct {
    symbols     []*strategy.Symbol
    interval    types.Interval
    strategy    strategy.Strategy
    datetime    *strategy.VnTime
    _dts        mapset.Set
    historyData map[string]map[strategy.VnTime]strategy.Bar
    *accountEngine
    gI gateway.GatewayInterface
}

var _LiveTradeEngine *LiveTradeEngine

func NewLiveTradeEngine(gI gateway.GatewayInterface) *LiveTradeEngine {
    if _LiveTradeEngine != nil {
        return _LiveTradeEngine
    }
    _LiveTradeEngine = &LiveTradeEngine{
        _dts:          mapset.NewSet(),
        accountEngine: newOrderEngine(),
        gI:            gI,
        historyData:   make(map[string]map[strategy.VnTime]strategy.Bar),
    }
    return _LiveTradeEngine
}

func (b *LiveTradeEngine) SetParameters(
    symbols []*strategy.Symbol,
    interval types.Interval,
) {
    b.symbols = symbols
    b.interval = interval
}

func (b *LiveTradeEngine) AddStrategy(strategy strategy.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    strategy.Inject(b.accountEngine)
    b.strategy = strategy
}

func (b *LiveTradeEngine) LoadData() {
    defer utils.TimeCost("load data")()
    for _, symbol := range b.symbols {
        bars, err := b.gI.LoadBarData(symbol)
        if err != nil {
            log.Fatal(err)
            panic(err)
        }
        if bars != nil {
            if b.historyData[symbol.Name] == nil {
                b.historyData[symbol.Name] = make(map[strategy.VnTime]strategy.Bar)
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
    dts := make([]strategy.VnTime, b._dts.Cardinality())
    cnt := 0
    b._dts.Each(func(ele interface{}) bool {
        dts[cnt] = ele.(strategy.VnTime)
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
    b.gI.WebSocketKLine(b.symbols)
    // b.vngo.
    // dts = make([]strategy.VnTime, b._dts.Cardinality())
    // cnt := 0
    // b._dts.Each(func(ele interface{}) bool {
    // 	dts[cnt] = ele.(strategy.VnTime)
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

func (b *LiveTradeEngine) newBars(dt strategy.VnTime) {
    bars := make(map[string]strategy.Bar)
    for _, symbol := range b.symbols {
        bars[symbol.Name] = b.historyData[symbol.Name][dt]
    }
    b.strategy.OnBars(bars)
}
