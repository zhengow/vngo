package backtesting_engine

import (
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/enum"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
    "strconv"
)

type backtestingAccount struct {
    engine.BaseAccount
    limitOrderCount int
    closes          map[strategy.Symbol]float64
}

func (o *backtestingAccount) Buy(symbol strategy.Symbol, price, volume float64) string {
    return o.sendOrder(symbol, enum.DirectionEnum.LONG, price, volume)
}

func (o *backtestingAccount) Sell(symbol strategy.Symbol, price, volume float64) string {
    return o.sendOrder(symbol, enum.DirectionEnum.SHORT, price, volume)
}

func (o *backtestingAccount) sendOrder(symbol strategy.Symbol, direction types.Direction, price, volume float64) string {
    price = o.PriceToTickSize(symbol, price)
    volume = o.VolumeToTickSize(symbol, volume)
    o.limitOrderCount++
    orderId := strconv.Itoa(o.limitOrderCount)
    order := strategy.NewOrder(symbol, orderId, direction, price, volume)
    o.Orders[orderId] = order
    return orderId
}

func (o *backtestingAccount) CancelAll() {
    o.Orders = make(map[string]*strategy.Order)
}

func (o *backtestingAccount) CancelById(orderId string) {
    delete(o.Orders, orderId)
}

func (o *backtestingAccount) startTrading() {
    //o.trading = true
}

func (o *backtestingAccount) GetPositions() map[strategy.Symbol]float64 {
    return o.Positions
}

func (o *backtestingAccount) updatePositions(symbol strategy.Symbol, incrementPos, price float64) {
    o.Positions[symbol] += incrementPos
    o.Cash -= incrementPos * price
}

//
//func (o *backtestingAccount) setPriceTicks(priceTicks map[string]int) {
//    o.priceTicks = priceTicks
//}

func (o *backtestingAccount) GetCash() float64 {
    return o.Cash
}

func (o *backtestingAccount) AddCash(increment float64) {
    o.Cash += increment
}

func (o *backtestingAccount) GetBalance() float64 {
    balance := o.Cash
    for symbol, position := range o.Positions {
        closePrice := o.closes[symbol]
        balance += closePrice * position
    }
    return balance
}

func (o *backtestingAccount) updateCloses(bars map[string]strategy.Bar) {
    for _, bar := range bars {
        o.closes[bar.Symbol] = bar.ClosePrice
    }
}

func newAccount() *backtestingAccount {
    account := new(backtestingAccount)
    account.Orders = make(map[string]*strategy.Order)
    account.Positions = make(map[strategy.Symbol]float64)
    account.closes = make(map[strategy.Symbol]float64)
    return account
}
