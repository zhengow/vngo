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

type BaseStrategy struct {
    accountInterface
    q *queue.Queue
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

func (s *BaseStrategy) OnInit(map[string]models.Bar) {
    s.q.Init.CloseChan()
}

func (s *BaseStrategy) OnBars(map[string]models.Bar) {
    s.q.Bars.CloseChan()
}

func (s *BaseStrategy) UpdateTrade(models.TradeData) {
    s.q.Trade.CloseChan()
}

func (s *BaseStrategy) UpdateOrder(*models.Order) {
    s.q.Order.CloseChan()
}

func (s *BaseStrategy) SetQueue(q *queue.Queue) {
    s.q = q
}

type Strategy interface {
    Inject(accountInterface)
    SetSetting(interface{}, map[string]interface{})
    OnInit(map[string]models.Bar)
    OnBars(map[string]models.Bar)
    UpdateTrade(models.TradeData)
    UpdateOrder(*models.Order)
    SetQueue(q *queue.Queue)
}
