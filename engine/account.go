package engine

import (
    "github.com/zhengow/vngo/enum"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/types"
    "github.com/zhengow/vngo/utils"
)

type accountEngine struct {
    priceTicks        map[string]int
    cash              float64
    activeLimitOrders map[int]*models.Order
    limitOrderCount   int
    closes            map[models.Symbol]float64
    positions         map[models.Symbol]float64
    //*positionEngine
}

func (o *accountEngine) Buy(symbol models.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, enum.DirectionEnum.LONG, price, volume)
}

func (o *accountEngine) Sell(symbol models.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, enum.DirectionEnum.SHORT, price, volume)
}

func (o *accountEngine) sendOrder(symbol models.Symbol, direction types.Direction, price, volume float64) int {
    priceTick := 5
    if val, ok := o.priceTicks[symbol.Name]; ok {
        priceTick = val
    }
    price = utils.RoundTo(price, priceTick)
    o.limitOrderCount++
    order := models.NewOrder(symbol, o.limitOrderCount, direction, price, volume)
    o.activeLimitOrders[o.limitOrderCount] = order
    return o.limitOrderCount
}

func (o *accountEngine) CancelAll() {
    o.activeLimitOrders = make(map[int]*models.Order)
}

func (o *accountEngine) startTrading() {
    //o.trading = true
}

func (o *accountEngine) GetPositions() map[models.Symbol]float64 {
    return o.positions
}

func (o *accountEngine) updatePositions(symbol models.Symbol, incrementPos, price float64) {
    o.positions[symbol] += incrementPos
    o.cash -= incrementPos * price
}

func (o *accountEngine) setPriceTicks(priceTicks map[string]int) {
    o.priceTicks = priceTicks
}

func (o *accountEngine) GetCash() float64 {
    return o.cash
}

func (o *accountEngine) AddCash(increment float64) {
    o.cash += increment
}

func (o *accountEngine) GetBalance() float64 {
    balance := o.cash
    for symbol, position := range o.positions {
        closePrice := o.closes[symbol]
        balance += closePrice * position
    }
    return balance
}

func (o *accountEngine) updateCloses(bars map[string]models.Bar) {
    for _, bar := range bars {
        o.closes[bar.Symbol] = bar.ClosePrice
    }
}

func newOrderEngine() *accountEngine {
    return &accountEngine{
        activeLimitOrders: make(map[int]*models.Order),
        positions:         make(map[models.Symbol]float64),
        closes:            make(map[models.Symbol]float64),
    }
}
