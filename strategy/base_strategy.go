package strategy

import (
    "fmt"
    "github.com/zhengow/vngo/model"
    "reflect"
    "time"
)

type accountInterface interface {
    Buy(symbol model.Symbol, price, volume float64) int
    Sell(symbol model.Symbol, price, volume float64) int
    CancelAll()
    GetPositions() map[model.Symbol]float64
    GetCash() float64
    GetBalance() float64
}

type marketInterface interface {
    CurrentTime() time.Time
}

type BaseStrategy struct {
    accountInterface
    marketInterface
    activeOrderIds []int
}

func (s *BaseStrategy) Inject(aI accountInterface) {
    s.accountInterface = aI
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
    Inject(accountInterface)
    SetSetting(interface{}, map[string]interface{})
    OnBars(map[string]model.Bar)
    UpdateTrade(model.TradeData)
}
