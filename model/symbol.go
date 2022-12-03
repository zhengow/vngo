package model

import "github.com/zhengow/vngo/consts"

type Symbol struct {
    Symbol   string
    Exchange consts.Exchange
}

func NewSymbol(symbol string, exchange consts.Exchange) *Symbol {
    return &Symbol{
        Symbol:   symbol,
        Exchange: exchange,
    }
}
