package strategy

import "github.com/zhengow/vngo/types"

type Symbol struct {
    Name     string `gorm:"column:symbol"`
    Exchange types.Exchange
    Interval types.Interval
    rate     float64
}

func (s *Symbol) Rate() float64 {
    return s.rate
}

func (s *Symbol) SetRate(rate float64) {
    s.rate = rate
}
