package model

import (
    "fmt"

    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/utils"
)

type Bar struct {
    Symbol
    Datetime     utils.DatabaseTime `gorm:"type:datetime"`
    Interval     consts.Interval
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
