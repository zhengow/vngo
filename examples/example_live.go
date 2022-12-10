package main

//import (
//   _ "embed"
//   "github.com/zhengow/vngo"
//   "log"
//
//   "github.com/zhengow/vngo/config"
//
//   "github.com/zhengow/vngo/database"
//   "github.com/zhengow/vngo/gateway"
//   "github.com/zhengow/vngo/live_trade_engine"
//)
//
////go:embed dev.yml
//var content []byte
//
//func getSymbols(symbols []string, exchange vngo.Exchange, interval vngo.Interval) []*vngo.Symbol {
//   res := make([]*vngo.Symbol, 0)
//   for _, symbol := range symbols {
//       res = append(res, vngo.NewSymbol(symbol, exchange, interval))
//   }
//   return res
//}
//
//func main() {
//   log.SetFlags(log.Ldate | log.Ltime)
//   _config, _ := config.NewConfig(content)
//   database.NewMysql(_config.MysqlConfig)
//   client := gateway.NewFutureClient(_config.Apikey, _config.SecretKey)
//   b := live_trade_engine.NewEngine(client)
//   symbols := getSymbols([]string{"BTCDOMUSDT"}, vngo.ExchangeEnum.BINANCE, vngo.IntervalEnum.MINUTE)
//   b.SetParameters(symbols, vngo.IntervalEnum.MINUTE)
//   b.AddStrategy(&MyStrategy{}, nil)
//   b.LoadData()
//   b.Run()
//}

func test() {}
