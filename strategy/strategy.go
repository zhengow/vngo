package strategy

import "reflect"

type Strategy struct {
	Test int
}

func (s *Strategy) SetSetting(setting map[string]interface{}) {
	fields := reflect.ValueOf(s).Elem()
	for name, value := range setting {
		filedValue := reflect.ValueOf(value)
		fields.FieldByName(name).Set(filedValue)
	}
}

