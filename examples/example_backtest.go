package main

import (
    _ "embed"
    "github.com/zhengow/vngo"
    "log"
    "time"
)

//go:embed dev.yml
var content []byte

func main() {
    log.SetFlags(log.Ldate | log.Ltime)
    _config, _ := vngo.NewConfig(content)
    vngo.UseMysql(_config.MysqlConfig)
    engine := vngo.NewBacktestingEngine().AddSymbols([]string{"BTCDOMUSDT"}, []float64{0.0001}, vngo.BinanceExchange).SetInterval(vngo.MinuteInterval)
    startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local)
    endDate := time.Date(2022, 7, 2, 0, 0, 0, 0, time.Local)
    engine.StartDate(startDate).EndDate(endDate).Capital(10000)
    engine.AddStrategy(&MyStrategy{Depth: 2}, nil)
    engine.LoadData()
    engine.RunBacktesting()
    engine.CalculateResult(true)
    //engine.ShowPNLChart()
    // engine.ShowKLineChart()
}
