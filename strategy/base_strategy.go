package strategy

import (
	"fmt"
	"reflect"
	"time"

	"github.com/zhengow/vngo/model"
)

type accountInterface interface {
	Buy(symbol model.Symbol, price, volume float64) int
	Sell(symbol model.Symbol, price, volume float64) int
	CancelAll()
	GetPositions() map[model.Symbol]float64
	CurrentTime() time.Time
    GetCash() float64
    GetBalance() float64
}

type BaseStrategy struct {
	accountInterface
	activeOrderIds []int
}

func (s *BaseStrategy) Inject(sI accountInterface) {
	s.accountInterface = sI
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
