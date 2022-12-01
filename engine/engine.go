package engine

import (
	"fmt"
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
}

var _BacktestingEngine *BacktestingEngine

func NewBacktestingEngine() *BacktestingEngine {
	if _BacktestingEngine != nil {
		return _BacktestingEngine
	}
	_BacktestingEngine = &BacktestingEngine{
		dts:         mapset.NewSet(),
		historyData: make(map[string]map[time.Time]model.Bar),
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
		fmt.Println("please set start and end time")
		return
	}
	start := b.start.Format(consts.DateFormat)
	end := b.end.Format(consts.DateFormat)
	for _, _symbol := range b.symbols {
		symbol, exchange := utils.ParseSymbol(_symbol)
		if symbol == "" || exchange == "" {
			continue
		}
		if b.historyData[symbol] == nil {
			b.historyData[symbol] = make(map[time.Time]model.Bar)
		}
		bars := database.LoadBarData(symbol, consts.Exchange(exchange), b.interval, start, end)
		for _, bar := range bars {
			_time := time.Time(bar.Datetime)
			b.dts.Add(_time)
			b.historyData[symbol][_time] = bar
		}
		fmt.Printf("%s load success, length: %d\n", symbol, len(b.historyData[symbol]))
	}

	// tmp := b.dts.ToSlice()
	// sort.Slice(tmp, func(i, j int) bool {
	// 	return tmp[i].(time.Time).Before(tmp[j].(time.Time))
	// })
}

func (b *BacktestingEngine) RunBacktesting() {
	b.strategy.Inject(newOrderEngine(b.priceTicks))
	b.strategy.OnInit()
	dts := make([]time.Time, b.dts.Cardinality())
	b.dts.Each(func(ele interface{}) bool {
		dts = append(dts, ele.(time.Time))
		return false
	})
	sort.Slice(dts, func(i, j int) bool {
		return dts[i].Before(dts[j])
	})

	// day_count := 0
	// idx := 0

	for _, dt := range dts {
		// dt := dts[idx]
		if b.datetime != nil && dt.Day() != b.datetime.Day() {
			b.strategy.DoneInit()
		}
		b.newBars(dt)
	}
	// println(b.strategy.Foo())
}

func (b *BacktestingEngine) newBars(dt time.Time) {
	b.datetime = &dt
	bars := make(map[string]model.Bar)
	for _, symbol := range b.symbols {
		bars[symbol] = b.historyData[symbol][dt]
	}
	b.crossLimitOrder()
	b.strategy.OnBars(bars)
}

func (b *BacktestingEngine) crossLimitOrder() {
	// println("cross limit order")
}

func (b *BacktestingEngine) CalculateResult() {
	println("calc")
}
