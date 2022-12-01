package engine

import (
	"github.com/zhengow/vngo/consts"
)

type orderEngine struct {
	priceTicks      map[string]int
	trading         bool
	activeOrderIds  []int
	limitOrderCount int
}

func (o *orderEngine) Buy(vtSymbol string, price, volume float64) *int {
	return o.sendOrder(vtSymbol, consts.DirectionEnum.LONG, price, volume)
}

func (o *orderEngine) Sell(vtSymbol string, price, volume float64) *int {
	return o.sendOrder(vtSymbol, consts.DirectionEnum.SHORT, price, volume)
}

func (o *orderEngine) sendOrder(vtSymbol string, direction consts.Direction, price, volume float64) *int {
	if !o.trading {
		return nil
	}
	// priceTick := 5
	// if val, ok := o.priceTicks[vtSymbol]; ok {
	// 	priceTick = val
	// }
	// price = utils.RoundTo(price, priceTick)
	// symbol, exchange := utils.ParseSymbol(vtSymbol)
	o.limitOrderCount++
	o.activeOrderIds = append(o.activeOrderIds, o.limitOrderCount)
	return &o.limitOrderCount
}

func (o *orderEngine) startTrading() {
	o.trading = true
}

func newOrderEngine(priceTicks map[string]int) *orderEngine {
	return &orderEngine{
		activeOrderIds: make([]int, 0),
		priceTicks:     priceTicks,
	}
}
