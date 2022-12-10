package vngo

import (
    "fmt"
    "github.com/zhengow/vngo/models"
    "reflect"
)

type accountInterface interface {
    Buy(symbol models.Symbol, price, volume float64) int
    Sell(symbol models.Symbol, price, volume float64) int
    CancelAll()
    GetPositions() map[models.Symbol]float64
    GetCash() float64
    GetBalance() float64
}

type marketInterface interface {
    CurrentTime() models.VnTime
}

type BaseStrategy struct {
    accountInterface
    marketInterface
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

func (s *BaseStrategy) OnBars(bars map[string]models.Bar) {
}

func (s *BaseStrategy) UpdateTrade(trade models.TradeData) {
}

type Strategy interface {
    Inject(accountInterface)
    SetSetting(interface{}, map[string]interface{})
    OnBars(map[string]models.Bar)
    UpdateTrade(models.TradeData)
}
