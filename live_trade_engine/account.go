package live_trade_engine

import (
    "fmt"

    "github.com/zhengow/vngo/model"
    "github.com/zhengow/vngo/utils"
)

type accountEngine struct {
    priceTicks        map[string]int
    cash              float64
    activeLimitOrders map[int]*model.Order
    limitOrderCount   int
    closes            map[model.Symbol]float64
    positions         map[model.Symbol]float64
    start             bool
}

func (o *accountEngine) Buy(symbol model.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, vngo.DirectionEnum.LONG, price, volume)
}

func (o *accountEngine) Sell(symbol model.Symbol, price, volume float64) int {
    return o.sendOrder(symbol, vngo.DirectionEnum.SHORT, price, volume)
}

func (o *accountEngine) sendOrder(symbol model.Symbol, direction vngo.Direction, price, volume float64) int {
    if !o.start {
        fmt.Println("trade")
        return -1
    }
    priceTick := 5
    if val, ok := o.priceTicks[symbol.Symbol]; ok {
        priceTick = val
    }
    price = utils.RoundTo(price, priceTick)
    o.limitOrderCount++
    order := model.NewOrder(symbol, o.limitOrderCount, direction, price, volume)
    o.activeLimitOrders[o.limitOrderCount] = order
    return o.limitOrderCount
}

func (o *accountEngine) CancelAll() {
    o.activeLimitOrders = make(map[int]*model.Order)
}

func (o *accountEngine) GetPositions() map[model.Symbol]float64 {
    return o.positions
}

func (o *accountEngine) updatePositions(symbol model.Symbol, incrementPos, price float64) {
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

func (o *accountEngine) updateCloses(bars map[string]model.Bar) {
    for _, bar := range bars {
        o.closes[bar.Symbol] = bar.ClosePrice
    }
}

func (o *accountEngine) startTrade() {
    o.start = true
}

func newOrderEngine() *accountEngine {
    return &accountEngine{
        activeLimitOrders: make(map[int]*model.Order),
        positions:         make(map[model.Symbol]float64),
        closes:            make(map[model.Symbol]float64),
    }
}
