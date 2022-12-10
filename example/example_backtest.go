package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/zhengow/vngo/config"
	"github.com/zhengow/vngo/model"

	"github.com/zhengow/vngo/backtesting_engine"
	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/database"
)

//go:embed dev.yml
var content []byte

func getSymbols(symbols []string, exchange consts.Exchange) []*model.Symbol {
	res := make([]*model.Symbol, 0)
	for _, symbol := range symbols {
		res = append(res, model.NewSymbol(symbol, exchange))
	}
	return res
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	b := backtesting_engine.NewEngine()
	symbols := getSymbols([]string{"BTCDOMUSDT"}, consts.ExchangeEnum.BINANCE)
	startDate := time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2022, 7, 2, 0, 0, 0, 0, time.Local)
	b.SetParameters(symbols, consts.IntervalEnum.MINUTE, startDate, endDate, nil, nil, 10000)
	b.AddStrategy(&MyStrategy{}, nil)
	_config, _ := config.NewConfig(content)
	database.NewMysql(_config.MysqlConfig)
	b.LoadData()
	b.RunBacktesting()
	b.CalculateResult(true)
	b.ShowPNLChart()
	// b.ShowKLineChart()
}
