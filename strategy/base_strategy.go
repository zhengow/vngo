package strategy

import (
	"github.com/zhengow/vngo/model"
)

type strategyInterface interface {
	Buy(vtSymbol string, price, volume float64) *int
	Sell(vtSymbol string, price, volume float64) *int
}

type BaseStrategy struct {
	strategyInterface
	init           bool
	trading        bool
	activeOrderIds []int
}

func (s *BaseStrategy) DoneInit() {
	s.init = true
}

// func (s *VirtualStrategy) Foo() bool {
// 	return s.init
// }

func (s *BaseStrategy) Inject(sI strategyInterface) {
	s.strategyInterface = sI
}

func (s *BaseStrategy) SetSetting(setting map[string]interface{}) {
	println("not implement set setting")
}

func (s *BaseStrategy) OnInit() {
	println("not implement on init")
}

func (s *BaseStrategy) OnStart() {
	println("not implement on start")
}

func (s *BaseStrategy) OnBars(bars map[string]model.Bar) {
	// println("not implement on bars")
}

type Strategy interface {
	Inject(strategyInterface)
	SetSetting(map[string]interface{})
	OnInit()
	DoneInit()
	OnStart()
	OnBars(map[string]model.Bar)
	// Foo() bool
}
