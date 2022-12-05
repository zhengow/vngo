package strategy

import (
    "fmt"
    "github.com/zhengow/vngo/model"
    "reflect"
)

type strategyInterface interface {
    Buy(symbol model.Symbol, price, volume float64) int
    Sell(symbol model.Symbol, price, volume float64) int
    CancelAll()
    GetPositions() map[model.Symbol]float64
}

type BaseStrategy struct {
    strategyInterface
    activeOrderIds []int
}

func (s *BaseStrategy) Inject(sI strategyInterface) {
    s.strategyInterface = sI
}

func (s *BaseStrategy) SetSetting(strategy interface{}, setting map[string]interface{}) {
    fields := reflect.ValueOf(strategy).Elem()
    for name, value := range setting {
        filedValue := reflect.ValueOf(value)
        if field := fields.FieldByName(name); field.CanSet() {
            field.Set(filedValue)
        } else {
            fmt.Printf("%s can't set\n", name)
        }
    }
}

func (s *BaseStrategy) OnBars(bars map[string]model.Bar) {
}

func (s *BaseStrategy) UpdateTrade(trade model.TradeData) {
}

type Strategy interface {
    Inject(strategyInterface)
    SetSetting(interface{}, map[string]interface{})
    OnBars(map[string]model.Bar)
    UpdateTrade(model.TradeData)
}
