package main

import (
    _ "embed"
    "github.com/zhengow/vngo"
    "github.com/zhengow/vngo/engine"
    "log"
    "time"
)

//go:embed dev.yml
var content []byte

func main() {
    log.SetFlags(log.Ldate | log.Ltime)
    b := engine.NewBacktestingEngine()
    symbols := vngo.GetSymbols([]string{"BTCDOMUSDT"}, vngo.BinanceExchange, vngo.MinuteInterval)
    startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local)
    endDate := time.Date(2022, 7, 2, 0, 0, 0, 0, time.Local)
    b.SetParameters(symbols, vngo.MinuteInterval, startDate, endDate, nil, nil, 10000)
    b.AddStrategy(&MyStrategy{Depth: 2}, nil)
    _config, _ := vngo.NewConfig(content)
    vngo.UseMysql(_config.MysqlConfig)
    b.LoadData()
    b.RunBacktesting()
    b.CalculateResult(true)
    b.ShowPNLChart()
    // b.ShowKLineChart()
}
