package engine

import (
    "fmt"
    "github.com/zhengow/vngo/enum"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/utils"
    "log"
    "math"
    "time"
)

type statisticEngine struct {
    capital  float64
    dts      []models.VnTime
    trades   map[int]*models.TradeData
    closes   map[models.VnTime]map[string]float64
    balances []float64
    rates    map[models.Symbol]float64
}

func (s *statisticEngine) setCapital(capital float64) {
    s.capital = capital
}

func (s *statisticEngine) updateClose(bars map[string]models.Bar) {
    currentCloses := make(map[string]float64)
    currentTime := models.VnTime{}
    for symbol, bar := range bars {
        currentCloses[symbol] = bar.ClosePrice
        currentTime = bar.Datetime
    }
    s.closes[currentTime] = currentCloses
    s.dts = append(s.dts, currentTime)
}

func getTrades(_trades map[int]*models.TradeData) map[models.VnTime][]*models.TradeData {
    trades := make(map[models.VnTime][]*models.TradeData)
    for _, trade := range _trades {
        if dtTrades, ok := trades[trade.Datetime]; ok {
            dtTrades = append(dtTrades, trade)
        } else {
            trades[trade.Datetime] = []*models.TradeData{trade}
        }
    }
    return trades
}

func (s *statisticEngine) CalculateResult(output bool) {
    if len(s.dts) == 0 {
        fmt.Println("no trade")
        return
    }
    trades := getTrades(s.trades)
    currentPos := make(map[string]float64)

    startDate, endDate := *models.NewVnTime(time.Time{}), s.dts[len(s.dts)-1]
    maxPNL, maxDrawdown, maxDrawdownPercent := s.capital, 0.0, 0.0
    totalTurnover, totalCommission := 0.0, 0.0
    totalTradeCount := 0
    dailyReturnStd, sharpeRatio := 0.0, 0.0

    retPercents := make([]float64, len(s.dts))
    retPercents[0] = 0
    balances := make([]float64, len(s.dts))
    balances[0] = s.capital
    for idx, dt := range s.dts[1:] {
        preCloses := s.closes[s.dts[idx]]
        closes := s.closes[dt]
        balance := balances[idx]
        for symbol, _close := range closes {
            pos := currentPos[symbol]
            if pos == 0 {
                continue
            }
            balance += pos * (_close - preCloses[symbol])
        }
        if dtTrades, ok := trades[dt]; ok {
            if startDate.IsZero() {
                startDate = dt
            }
            for _, _trade := range dtTrades {
                symbol := _trade.Symbol.Name
                volume := _trade.Volume
                if _trade.Direction == enum.DirectionEnum.SHORT {
                    volume *= -1
                }
                currentPos[symbol] += volume
                rate := 0.0
                if r, ok := s.rates[_trade.Symbol]; ok {
                    rate = r
                }
                turnover := _trade.Price * _trade.Volume
                fee := rate * _trade.Price * _trade.Volume
                totalCommission += fee
                totalTurnover += turnover
                balance += volume*(closes[symbol]-_trade.Price) - fee
                totalTradeCount++
            }
        }
        balances[idx+1] = balance
        retPercents[idx+1] = balances[idx+1]/balances[idx] - 1
        maxPNL = math.Max(maxPNL, balance)
        maxDrawdown = math.Min(maxDrawdown, balance-maxPNL)
        maxDrawdownPercent = math.Min(maxDrawdownPercent, (balance-maxPNL)/maxPNL)
    }
    s.balances = balances

    if output {
        totalMinutes := endDate.Sub(startDate.Time).Minutes()
        totalRetAmount := s.balances[len(s.balances)-1] - s.balances[0]
        totalRetPercent := totalRetAmount / s.balances[0]
        annualRet := totalRetPercent / totalMinutes * utils.YearMinutes()
        dailyTradeCount := float64(totalTradeCount) / totalMinutes * (time.Hour.Minutes() * 24)
        dailyRetAmount := totalRetAmount / totalMinutes * (time.Hour.Minutes() * 24)
        dailyRetPercent := totalRetPercent / totalMinutes * (time.Hour.Minutes() * 24)
        //fmt.Println(retPercents)
        dailyReturnStd = utils.Std(retPercents) * math.Sqrt(time.Hour.Minutes()*24)
        sharpeRatio = dailyRetPercent / dailyReturnStd * math.Sqrt(365) // 算出每日夏普
        // returnDropdownRatio = -annualRet / maxDrawdownPercent
        log.Printf("%-15s\t%-s\n", "首个交易日：", startDate.Format())
        log.Printf("%-15s\t%-s\n", "最后交易日：", endDate.Format())
        log.Printf("%-15s\t%-.2f\n", "起始资金：", s.balances[0])
        log.Printf("%-15s\t%-.2f\n", "结束资金：", s.balances[len(s.balances)-1])
        log.Printf("%-15s\t%-.2f\n", "最高资金：", maxPNL)
        log.Printf("%-15s\t\t%-.2f\n", "总盈亏：", s.balances[len(s.balances)-1]-s.balances[0])
        log.Printf("%-15s\t%-.2f%%\n", "总收益率：", totalRetPercent*100)
        log.Printf("%-15s\t%-.2f%%\n", "年化收益率：", annualRet*100)
        log.Printf("%-15s\t%-.2f\n", "最大回撤：", maxDrawdown)
        log.Printf("%-15s\t%-.2f%%\n", "百分比最大回撤：", maxDrawdownPercent*100)
        log.Printf("%-15s\t%-.2f\n", "总手续费：", totalCommission)
        log.Printf("%-15s\t%-.2f\n", "总成交金额：", totalTurnover)
        log.Printf("%-15s\t%-d\n", "总成交笔数：", totalTradeCount)
        log.Printf("%-15s\t%-.2f\n", "日均盈亏：", dailyRetAmount)
        log.Printf("%-15s\t%-.2f\n", "日均成交笔数：", dailyTradeCount)
        log.Printf("%-15s\t%-.2f%%\n", "日均收益率：", dailyRetPercent*100)
        log.Printf("%-15s\t%-.2f%%\n", "日均标准差：", dailyReturnStd*100)
        log.Printf("%-15s\t%-.2f\n", "夏普比率：", sharpeRatio)
        // log.Printf("%-15s\t%-.2f\n", "收益回撤比：", returnDropdownRatio)
    }
}

func newStatisticEngine() *statisticEngine {
    return &statisticEngine{
        trades: make(map[int]*models.TradeData),
        closes: make(map[models.VnTime]map[string]float64),
    }
}

func (s *statisticEngine) setRates(rates map[models.Symbol]float64) {
    s.rates = rates
}

//
//type eachResult struct {
//    date models.VnTime
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
