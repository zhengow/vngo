package model

import (
    "github.com/zhengow/vngo/consts"
)

type Order struct {
    Symbol    Symbol
    OrderId   int
    Direction consts.Direction
    Price     float64
    Volume    float64
    //status    consts.Status
    //datetime time.Time
}

func NewOrder(symbol Symbol,
    orderId int,
    direction consts.Direction,
    price float64,
    volume float64,
//status consts.Status,
//    datetime time.Time
) *Order {
    return &Order{
        Symbol:    symbol,
        OrderId:   orderId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        //status:    status,
        //datetime: datetime,
    }
}
