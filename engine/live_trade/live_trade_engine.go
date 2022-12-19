package live_trade

import (
	"fmt"
	"log"
	"sort"

	mapset "github.com/deckarep/golang-set"
	"github.com/zhengow/vngo/engine"
	"github.com/zhengow/vngo/gateway"
	"github.com/zhengow/vngo/models"
	"github.com/zhengow/vngo/utils"
)

type Engine struct {
	*engine.BaseEngine
	datetime    *models.VnTime
	_dts        mapset.Set
	historyData map[string]map[models.VnTime]models.Bar
	gI          gateway.GatewayInterface
}

func NewLiveTradeEngine(gI gateway.GatewayInterface) *Engine {
	return &Engine{
		_dts:        mapset.NewSet(),
		gI:          gI,
		historyData: make(map[string]map[models.VnTime]models.Bar),
		BaseEngine:  engine.NewBaseEngine(models.EngineEnum.LIVETRADEENGINE),
	}
}

func (b *Engine) LoadHistoryData() {
	defer utils.TimeCost("load data")()
	for _, symbol := range b.GetSymbols() {
		bars, err := b.gI.LoadBarData(symbol, b.GetInterval())
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		if bars != nil {
			if b.historyData[symbol.FullName()] == nil {
				b.historyData[symbol.FullName()] = make(map[models.VnTime]models.Bar)
			}
			for _, bar := range bars {
				_time := bar.Datetime
				b._dts.Add(_time)
				b.historyData[symbol.FullName()][_time] = bar
			}
			fmt.Printf("%s load success, length: %d\n", symbol.FullName(), len(b.historyData[symbol.FullName()]))
		}
	}
	b.preRun()
}

func (b *Engine) preRun() {
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

func (b *Engine) Run() {
	if b.GetInterval() == models.IntervalEnum.MINUTE {
		b.gI.LoadBarDataByMinute(b.GetSymbols(), b.Queue)
	}
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

func (b *Engine) newBars(dt models.VnTime) {
	bars := make(map[string]models.Bar)
	for _, symbol := range b.GetSymbols() {
		bars[symbol.FullName()] = b.historyData[symbol.FullName()][dt]
	}
	b.Queue.Init.SendSync(bars)
}
