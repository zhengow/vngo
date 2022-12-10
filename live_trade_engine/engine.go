package live_trade_engine

import (
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/gateway"
	"github.com/zhengow/vngo/model"
	"github.com/zhengow/vngo/strategy"
	"github.com/zhengow/vngo/utils"
)

type LiveTradeEngine struct {
	symbols  []*model.Symbol
	interval consts.Interval
	strategy strategy.Strategy
	datetime *time.Time
	*accountEngine
	gI gateway.GatewayInterface
}

var _LiveTradeEngine *LiveTradeEngine

func NewEngine(gI gateway.GatewayInterface) *LiveTradeEngine {
	if _LiveTradeEngine != nil {
		return _LiveTradeEngine
	}
	_LiveTradeEngine = &LiveTradeEngine{
		accountEngine: newOrderEngine(),
		gI:            gI,
	}
	return _LiveTradeEngine
}

func (b *LiveTradeEngine) SetParameters(
	symbols []*model.Symbol,
	interval consts.Interval,
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
}

func (b *LiveTradeEngine) Run() {
	// b.dts = make([]time.Time, b._dts.Cardinality())
	// cnt := 0
	// b._dts.Each(func(ele interface{}) bool {
	// 	b.dts[cnt] = ele.(time.Time)
	// 	cnt++
	// 	return false
	// })
	// sort.Slice(b.dts, func(i, j int) bool {
	// 	return b.dts[i].Before(b.dts[j])
	// })

	// for _, dt := range b.dts {
	// 	b.newBars(dt)
	// }
}

func (b *LiveTradeEngine) newBars(dt time.Time) {
	b.datetime = &dt
	bars := make(map[string]model.Bar)
	// for _, symbol := range b.symbols {
	// 	bars[symbol.Symbol] = b.historyData[symbol.Symbol][dt]
	// }
	b.accountEngine.updateCloses(bars)
	b.strategy.OnBars(bars)
}
