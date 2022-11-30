package strategy

import "github.com/zhengow/vngo/model"

type VirtualStrategy struct {
	init bool
	trading bool
	activeOrderIds []int
}

func (s *VirtualStrategy) DoneInit() {
	s.init = true
}

// func (s *VirtualStrategy) Foo() bool {
// 	return s.init
// }

func (s *VirtualStrategy) SetSetting(setting map[string]interface{}) {
	println("not implement set setting")
}

func (s *VirtualStrategy) OnInit() {
	println("not implement on init")
}

func (s *VirtualStrategy) OnStart() {
	println("not implement on start")
}

func (s *VirtualStrategy) OnBars(bars map[string]model.Bar) {
	// println("not implement on bars")
}

func (s *VirtualStrategy) sendOrder(bars map[string]model.Bar) {
	// println("not implement on bars")
}

type Strategy interface {
	SetSetting(setting map[string]interface{})
	OnInit()
	DoneInit()
	OnStart()
	OnBars(map[string]model.Bar)
	// Foo() bool
}
