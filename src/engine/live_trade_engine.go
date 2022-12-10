package engine

import (
    "fmt"
    "github.com/zhengow/vngo"
    "log"
    "sort"
    "time"

    mapset "github.com/deckarep/golang-set"
    "github.com/zhengow/vngo/gateway"
)

type LiveTradeEngine struct {
    symbols     []*vngo.Symbol
    interval    vngo.Interval
    strategy    vngo.Strategy
    datetime    *time.Time
    _dts        mapset.Set
    historyData map[string]map[time.Time]vngo.Bar
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
        historyData:   make(map[string]map[time.Time]vngo.Bar),
    }
    return _LiveTradeEngine
}

func (b *LiveTradeEngine) SetParameters(
    symbols []*vngo.Symbol,
    interval vngo.Interval,
) {
    b.symbols = symbols
    b.interval = interval
}

func (b *LiveTradeEngine) AddStrategy(strategy vngo.Strategy, setting map[string]interface{}) {
    strategy.SetSetting(strategy, setting)
    strategy.Inject(b.accountEngine)
    b.strategy = strategy
}

func (b *LiveTradeEngine) LoadData() {
    defer vngo.TimeCost("load data")()
    for _, symbol := range b.symbols {
        bars, err := b.gI.LoadBarData(symbol)
        if err != nil {
            log.Fatal(err)
            panic(err)
        }
        if bars != nil {
            if b.historyData[symbol.Symbol] == nil {
                b.historyData[symbol.Symbol] = make(map[time.Time]vngo.Bar)
            }
            for _, bar := range bars {
                _time := bar.Datetime.Time
                b._dts.Add(_time)
                b.historyData[symbol.Symbol][_time] = bar
            }
            fmt.Printf("%s.%s load success, length: %d\n", symbol.Symbol, symbol.Exchange, len(b.historyData[symbol.Symbol]))
        }
    }
    b.preRun()
}

func (b *LiveTradeEngine) preRun() {
    dts := make([]time.Time, b._dts.Cardinality())
    cnt := 0
    b._dts.Each(func(ele interface{}) bool {
        dts[cnt] = ele.(time.Time)
        cnt++
        return false
    })
    sort.Slice(dts, func(i, j int) bool {
        return dts[i].Before(dts[j])
    })

    for _, dt := range dts {
        b.newBars(dt)
    }
}

func (b *LiveTradeEngine) Run() {
    b.gI.WebSocketKLine(b.symbols)
    // b.vngo.
    // dts = make([]time.Time, b._dts.Cardinality())
    // cnt := 0
    // b._dts.Each(func(ele interface{}) bool {
    // 	dts[cnt] = ele.(time.Time)
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

func (b *LiveTradeEngine) newBars(dt time.Time) {
    bars := make(map[string]vngo.Bar)
    for _, symbol := range b.symbols {
        bars[symbol.Symbol] = b.historyData[symbol.Symbol][dt]
    }
    b.strategy.OnBars(bars)
}
