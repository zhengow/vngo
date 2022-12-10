package main

import (
	_ "embed"
	"log"

	"github.com/zhengow/vngo/config"
	"github.com/zhengow/vngo/model"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/database"
	"github.com/zhengow/vngo/gateway"
	"github.com/zhengow/vngo/live_trade_engine"
)

//go:embed dev.yml
var content []byte

func getSymbols(symbols []string, exchange consts.Exchange, interval consts.Interval) []*model.Symbol {
	res := make([]*model.Symbol, 0)
	for _, symbol := range symbols {
		res = append(res, model.NewSymbol(symbol, exchange, interval))
	}
	return res
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	_config, _ := config.NewConfig(content)
	database.NewMysql(_config.MysqlConfig)
	client := gateway.NewFutureClient(_config.Apikey, _config.SecretKey)
	b := live_trade_engine.NewEngine(client)
	symbols := getSymbols([]string{"BTCDOMUSDT"}, consts.ExchangeEnum.BINANCE, consts.IntervalEnum.MINUTE)
	b.SetParameters(symbols, consts.IntervalEnum.MINUTE)
	b.AddStrategy(&MyStrategy{}, nil)
	b.LoadData()
	// b.Run()
}
