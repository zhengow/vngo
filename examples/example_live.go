package main

import (
    _ "embed"
    "github.com/zhengow/vngo"
    "github.com/zhengow/vngo/gateway"
    "log"
)

//go:embed dev.yml
var content []byte

func main() {
    log.SetFlags(log.Ldate | log.Ltime)
    _config, _ := vngo.NewConfig(content)
    vngo.UseMysql(_config.MysqlConfig)
    client := gateway.NewFutureClient(_config.Apikey, _config.SecretKey)
    b := vngo.NewLiveTradeEngine(client)
    b.AddSymbol("BTCDOMUSDT", vngo.BinanceExchange).SetInterval(vngo.MinuteInterval)
    b.AddStrategy(&MyStrategy{}, nil)
    b.LoadData()
    b.Run()
}
