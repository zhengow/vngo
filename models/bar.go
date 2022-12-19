package models

import (
    "fmt"
    "time"
)

type Bar struct {
    Symbol
    Datetime     VnTime `gorm:"type:datetime"`
    Interval     Interval
    Volume       float64
    OpenInterest float64
    OpenPrice    float64
    HighPrice    float64
    LowPrice     float64
    ClosePrice   float64
}

func (b *Bar) FullName() string {
    return fmt.Sprintf("%s.%s", b.Symbol.Name, b.Symbol.Exchange)
}

func (b *Bar) GetKLineData() [4]float64 {
    return [4]float64{b.OpenPrice, b.ClosePrice, b.LowPrice, b.HighPrice}
}

func (b *Bar) GetDatetime() time.Time {
    return b.Datetime.Time
}

//
//func (b *Bar) GetSymbol() Symbol {
//    return b.Symbol
//}
//
//func (b *Bar) GetInterval() models.Interval {
//    return b.Interval
//}
//
//func (b *Bar) GetVolume() float64 {
//    return b.Volume
//}
//
//func (b *Bar) GetOpenInterest() float64 {
//    return b.OpenInterest
//}
//
//func (b *Bar) GetOpenPrice() float64 {
//    return b.OpenPrice
//}
//
//func (b *Bar) GetHighPrice() float64 {
//    return b.HighPrice
//}
//
//func (b *Bar) GetLowPrice() float64 {
//    return b.LowPrice
//}
//
//func (b *Bar) GetClosePrice() float64 {
//    return b.ClosePrice
//}
