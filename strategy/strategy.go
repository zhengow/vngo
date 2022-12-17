package strategy

import (
    "fmt"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/queue"
    "reflect"
)

type accountInterface interface {
    Buy(symbol models.Symbol, price, volume float64) string
    Sell(symbol models.Symbol, price, volume float64) string
    CancelAll()
    CancelById(orderId string)
    GetPositions() map[models.Symbol]float64
    GetCash() float64
    GetBalance() float64
}

//type marketInterface interface {
//    CurrentTime() VnTime
//}

type BaseStrategy struct {
    accountInterface
    queue.Queue
    //marketInterface
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

func (s *BaseStrategy) OnBars(map[models.Symbol]models.Bar) {
}

func (s *BaseStrategy) UpdateTrade(models.TradeData) {
}

func (s *BaseStrategy) UpdateOrder(*models.Order) {
}

type Strategy interface {
    Inject(accountInterface)
    SetSetting(interface{}, map[string]interface{})
    OnBars(map[string]models.Bar)
    UpdateTrade(models.TradeData)
    UpdateOrder(*models.Order)
}

//type models.Bar interface {
//    GetDateTime() time.Time
//    Getmodels.Symbol() models.Symbol
//    GetOpen() float64
//    GetHigh() float64
//    GetLow() float64
//    GetClose() float64
//    GetVolume() float64
//}
//
//type models.Symbol interface {
//    GetName() string
//    GetExchange() string
//}
//
//type models.TradeData interface {
//    IsSell() bool
//    GetDateTime() time.Time
//    Getmodels.Symbol() models.Symbol
//    GetPrice() float64
//    GetVolume() float64
//}
