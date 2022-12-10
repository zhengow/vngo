package backtesting_engine

import (
    "github.com/zhengow/vngo"
)

type accountEngine struct {
    priceTicks        map[string]int
    cash              float64
    activeLimitOrders map[int]*vngo.Order
    limitOrderCount   int
    closes            map[vngo.Symbol]float64
    positions         map[vngo.Symbol]float64
    //*positionEngine
}

func (o *accountEngine) Buy(symbol vngo.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, vngo.DirectionEnum.LONG, price, volume)
}

func (o *accountEngine) Sell(symbol vngo.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, vngo.DirectionEnum.SHORT, price, volume)
}

func (o *accountEngine) sendOrder(symbol vngo.Symbol, direction vngo.Direction, price, volume float64) int {
    priceTick := 5
    if val, ok := o.priceTicks[symbol.Symbol]; ok {
        priceTick = val
    }
    price = vngo.RoundTo(price, priceTick)
    o.limitOrderCount++
    order := vngo.NewOrder(symbol, o.limitOrderCount, direction, price, volume)
    o.activeLimitOrders[o.limitOrderCount] = order
    return o.limitOrderCount
}

func (o *accountEngine) CancelAll() {
    o.activeLimitOrders = make(map[int]*vngo.Order)
}

func (o *accountEngine) startTrading() {
    //o.trading = true
}

func (o *accountEngine) GetPositions() map[vngo.Symbol]float64 {
    return o.positions
}

func (o *accountEngine) updatePositions(symbol vngo.Symbol, incrementPos, price float64) {
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

func (o *accountEngine) updateCloses(bars map[string]vngo.Bar) {
    for _, bar := range bars {
        o.closes[bar.Symbol] = bar.ClosePrice
    }
}

func newOrderEngine() *accountEngine {
    return &accountEngine{
        activeLimitOrders: make(map[int]*vngo.Order),
        positions:         make(map[vngo.Symbol]float64),
        closes:            make(map[vngo.Symbol]float64),
    }
}
