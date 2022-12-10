package vngo

import (
    "fmt"
)

type Bar struct {
    Symbol
    Datetime     DatabaseTime `gorm:"type:datetime"`
    Interval     Interval
    Volume       float64
    OpenInterest float64
    OpenPrice    float64
    HighPrice    float64
    LowPrice     float64
    ClosePrice   float64
}

func (b *Bar) FullName() string {
    return fmt.Sprintf("%s.%s", b.Symbol, b.Exchange)
}

func (b *Bar) GetKLineData() [4]float64 {
    return [4]float64{b.OpenPrice, b.ClosePrice, b.LowPrice, b.HighPrice}
}
