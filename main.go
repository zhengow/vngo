package main

import (
    "github.com/zhengow/vngo/model"
    "time"

    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/engine"
    "github.com/zhengow/vngo/strategy"
)

func getSymbols(symbols []string, exchange consts.Exchange) []*model.Symbol {
    res := make([]*model.Symbol, 0)
    for _, symbol := range symbols {
        res = append(res, model.NewSymbol(symbol, exchange))
    }
    return res
}

func main() {
    b := engine.NewBacktestingEngine()
    symbols := getSymbols([]string{"BTCDOMUSDT"}, consts.ExchangeEnum.BINANCE)
    b.SetParameters(symbols, consts.IntervalEnum.MINUTE, time.Now().AddDate(0, -1, 0), time.Now(), nil, 10000)
    b.AddStrategy(&strategy.MyStrategy{}, map[string]interface{}{"Test": 1})
    database.NewMysql()
    b.LoadData()
    b.RunBacktesting()
    b.CalculateResult()
}
