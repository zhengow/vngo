package strategy

import (
    "github.com/zhengow/vngo/model"
)

type strategyInterface interface {
    Buy(symbol model.Symbol, price, volume float64) int
    Sell(symbol model.Symbol, price, volume float64) int
}

type BaseStrategy struct {
    strategyInterface
    activeOrderIds []int
}

func (s *BaseStrategy) Inject(sI strategyInterface) {
    s.strategyInterface = sI
}

//func (s *BaseStrategy) IsInit() bool {
//    return s.init
//}

func (s *BaseStrategy) SetSetting(setting map[string]interface{}) {
    println("not implement set setting")
}

func (s *BaseStrategy) OnInit() {
    println("not implement on init")
}

func (s *BaseStrategy) DoneInit() {
    //s.init = true
}

func (s *BaseStrategy) OnStart() {
    println("not implement on start")
}

func (s *BaseStrategy) OnBars(bars map[string]model.Bar) {
    println("not implement on bars")
}

func (s *BaseStrategy) UpdateTrade(trade model.TradeData) {
    //println("not implement update trade")
}

type Strategy interface {
    Inject(strategyInterface)
    //IsInit() bool
    SetSetting(map[string]interface{})
    OnInit()
    DoneInit()
    OnStart()
    OnBars(map[string]model.Bar)
    UpdateTrade(model.TradeData)
}
