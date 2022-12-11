package strategy

import "github.com/zhengow/vngo/types"

type Symbol struct {
    Name     string `gorm:"column:symbol"`
    Exchange types.Exchange
}
