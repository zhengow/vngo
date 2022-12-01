package strategy

import "reflect"

type MyStrategy struct {
	BaseStrategy
	Test int
}

func (s MyStrategy) SetSetting(setting map[string]interface{}) {
	fields := reflect.ValueOf(&s).Elem()
	for name, value := range setting {
		filedValue := reflect.ValueOf(value)
		fields.FieldByName(name).Set(filedValue)
	}
}

func (s MyStrategy) OnInit() {
	println("on init")
}
