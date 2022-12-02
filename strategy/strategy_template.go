package strategy

import (
    "github.com/zhengow/vngo/model"
    "reflect"
)

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

func (s MyStrategy) OnBars(bars map[string]model.Bar) {
    for symbol, bar := range bars {
        //fmt.Println(symbol, bar.Datetime.Format())
        s.Buy(symbol, bar.ClosePrice, 1)
    }
}
