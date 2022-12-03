package engine

import (
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/model"
    "github.com/zhengow/vngo/utils"
)

type orderEngine struct {
    priceTicks map[string]int
    //trading           bool
    activeLimitOrders map[int]*model.Order
    limitOrderCount   int
}

func (o *orderEngine) Buy(symbol model.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, consts.DirectionEnum.LONG, price, volume)
}

func (o *orderEngine) Sell(symbol model.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, consts.DirectionEnum.SHORT, price, volume)
}

func (o *orderEngine) sendOrder(symbol model.Symbol, direction consts.Direction, price, volume float64) int {
    //if !o.trading {
    //    return -1
    //}
    priceTick := 5
    if val, ok := o.priceTicks[symbol.Symbol]; ok {
        priceTick = val
    }
    price = utils.RoundTo(price, priceTick)
    o.limitOrderCount++
    order := model.NewOrder(symbol.Symbol, symbol.Exchange, o.limitOrderCount, direction, price, volume)
    o.activeLimitOrders[o.limitOrderCount] = order
    return o.limitOrderCount
}

func (o *orderEngine) startTrading() {
    //o.trading = true
}

func newOrderEngine(priceTicks map[string]int) *orderEngine {
    return &orderEngine{
        activeLimitOrders: make(map[int]*model.Order),
        priceTicks:        priceTicks,
    }
}