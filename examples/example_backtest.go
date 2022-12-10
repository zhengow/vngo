package main

import (
    _ "embed"
    "github.com/zhengow/vngo"
    "log"
    "time"

    "github.com/zhengow/vngo/engine"
)

//go:embed dev.yml
var content []byte

func getSymbols(symbols []string, exchange vngo.Exchange, interval vngo.Interval) []*vngo.Symbol {
    res := make([]*vngo.Symbol, 0)
    for _, symbol := range symbols {
        res = append(res, vngo.NewSymbol(symbol, exchange, interval))
    }
    return res
}

func main() {
    log.SetFlags(log.Ldate | log.Ltime)
    b := engine.NewBacktestingEngine()
    symbols := getSymbols([]string{"BTCDOMUSDT"}, vngo.ExchangeEnum.BINANCE, vngo.IntervalEnum.MINUTE)
    startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local)
    endDate := time.Date(2022, 7, 2, 0, 0, 0, 0, time.Local)
    b.SetParameters(symbols, vngo.IntervalEnum.MINUTE, startDate, endDate, nil, nil, 10000)
    b.AddStrategy(&MyStrategy{Depth: 2}, nil)
    _config, _ := vngo.NewConfig(content)
    vngo.UseMysql(_config.MysqlConfig)
    b.LoadData()
    b.RunBacktesting()
    b.CalculateResult(true)
    b.ShowPNLChart()
    // b.ShowKLineChart()
}
