package vngo

import (
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/types"
)

func newSymbol(name string, exchange types.Exchange, interval types.Interval) *models.Symbol {
    return &models.Symbol{
        Name:     name,
        Exchange: exchange,
        Interval: interval,
    }
}

func GetSymbols(symbols []string, exchange types.Exchange, interval types.Interval) []*models.Symbol {
    res := make([]*models.Symbol, 0)
    for _, symbol := range symbols {
        res = append(res, newSymbol(symbol, exchange, interval))
    }
    return res
}
