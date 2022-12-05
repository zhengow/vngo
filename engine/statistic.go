package engine

import (
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/model"
)

type statisticEngine struct {
	capital float64
	dts     []time.Time
	trades  map[int]*model.TradeData
	closes  map[time.Time]map[string]float64
	netPnls []float64
	rates   map[model.Symbol]float64
}

func (s *statisticEngine) setCaptial(capital float64) {
	s.capital = capital
}

func (s *statisticEngine) updateClose(bars map[string]model.Bar) {
	currentCloses := make(map[string]float64)
	currentTime := time.Time{}
	for symbol, bar := range bars {
		currentCloses[symbol] = bar.ClosePrice
		currentTime = time.Time(bar.Datetime)
	}
	s.closes[currentTime] = currentCloses
	s.dts = append(s.dts, currentTime)
}

func (s *statisticEngine) CalculateResult() {
	trades := make(map[time.Time][]*model.TradeData)
	for _, trade := range s.trades {
		if dtTrades, ok := trades[trade.Datetime]; ok {
			dtTrades = append(dtTrades, trade)
		} else {
			trades[trade.Datetime] = []*model.TradeData{trade}
		}
	}
	currentPos := make(map[string]float64)
	netPnls := make([]float64, len(s.dts))
	netPnls[0] = s.capital
	for idx, dt := range s.dts[1:] {
		preCloses := s.closes[s.dts[idx]]
		closes := s.closes[dt]
		pnl := netPnls[idx]
		for symbol, _close := range closes {
			pos := currentPos[symbol]
			if pos == 0 {
				continue
			}
			pnl += pos * (_close - preCloses[symbol])
		}
		if dtTrades, ok := trades[dt]; ok {
			for _, _trade := range dtTrades {
				symbol := _trade.Symbol.Symbol
				volume := _trade.Volume
				if _trade.Direction == consts.DirectionEnum.SHORT {
					volume *= -1
				}
				currentPos[symbol] += volume
				rate := 0.0
				if r, ok := s.rates[_trade.Symbol]; ok {
					rate = r
				}
				pnl += volume * (closes[symbol] - _trade.Price) - rate * _trade.Price * _trade.Volume
			}
		}
		netPnls[idx+1] = pnl
	}
	s.netPnls = netPnls
}

func newStatisticEngine() *statisticEngine {
	return &statisticEngine{
		trades: make(map[int]*model.TradeData),
		closes: make(map[time.Time]map[string]float64),
	}
}

func (s *statisticEngine) setRates(rates map[model.Symbol]float64) {
	s.rates = rates
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
