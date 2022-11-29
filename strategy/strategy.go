package strategy

type VirtualStrategy struct {
	Init bool
}

func (s VirtualStrategy) DoneInit() {
	s.Init = true
	_ = s.Init
}

func (s VirtualStrategy) SetSetting(setting map[string]interface{}) {
	println("not implement set setting")
}

func (s VirtualStrategy) OnInit() {
	println("not implement on init")
}

func (s VirtualStrategy) OnStart() {
	println("not implement on start")
}

type Strategy interface {
	SetSetting(setting map[string]interface{})
	OnInit()
	DoneInit()
	OnStart()
}
