package engine

import (
    "github.com/zhengow/vngo/model"
    "time"
)

type statistic struct {
    trades map[int]*model.TradeData
    closes map[time.Time]map[string]float64
}

//
//type eachResult struct {
//    date time.Time
//    closePrices map[string]float64
//    preCloses map[string]float64
//    startPoses map[string]float64
//    endPoses map[string]float64
//    contract_results Dict[str, ContractDailyResult] = {}
//
//    for vt_symbol, closePrice in closePrices.items()
//    contract_results[vt_symbol] = ContractDailyResult(result_date, closePrice)
//
//    tradeCount int = 0
//    turnover float = 0
//    commission float = 0
//    slippage float = 0
//    tradingPnl float = 0
//    holdingPnl float = 0
//    totalPnl float = 0
//    netPnl float = 0
//}
