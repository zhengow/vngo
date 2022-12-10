package vngo

import (
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
)

func newSymbol(name string, exchange types.Exchange, interval types.Interval) *strategy.Symbol {
    return &strategy.Symbol{
        Name:     name,
        Exchange: exchange,
        Interval: interval,
    }
}

func GetSymbols(symbols []string, exchange types.Exchange, interval types.Interval) []*strategy.Symbol {
    res := make([]*strategy.Symbol, 0)
    for _, symbol := range symbols {
        res = append(res, newSymbol(symbol, exchange, interval))
    }
    return res
}
