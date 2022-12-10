package strategy

import (
    "fmt"
    "reflect"
)

type accountInterface interface {
    Buy(symbol Symbol, price, volume float64) string
    Sell(symbol Symbol, price, volume float64) string
    CancelAll()
    CancelById(orderId string)
    GetPositions() map[Symbol]float64
    GetCash() float64
    GetBalance() float64
}

type marketInterface interface {
    CurrentTime() VnTime
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

func (s *BaseStrategy) OnBars(bars map[string]Bar) {
}

func (s *BaseStrategy) UpdateTrade(trade TradeData) {
}

type Strategy interface {
    Inject(accountInterface)
    SetSetting(interface{}, map[string]interface{})
    OnBars(map[string]Bar)
    UpdateTrade(TradeData)
}

//type Bar interface {
//    GetDateTime() time.Time
//    GetSymbol() Symbol
//    GetOpen() float64
//    GetHigh() float64
//    GetLow() float64
//    GetClose() float64
//    GetVolume() float64
//}
//
//type Symbol interface {
//    GetName() string
//    GetExchange() string
//}
//
//type TradeData interface {
//    IsSell() bool
//    GetDateTime() time.Time
//    GetSymbol() Symbol
//    GetPrice() float64
//    GetVolume() float64
//}
