package model

import "github.com/zhengow/vngo/consts"

type Symbol struct {
    Symbol   string
    Exchange consts.Exchange
    Interval consts.Interval
}

func NewSymbol(symbol string, exchange consts.Exchange, interval consts.Interval) *Symbol {
    return &Symbol{
        Symbol:   symbol,
        Exchange: exchange,
        Interval: interval,
    }
}
