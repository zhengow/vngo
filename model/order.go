package model

import (
    "github.com/zhengow/vngo/consts"
)

type Order struct {
    Symbol    string
    Exchange  consts.Exchange
    OrderId   int
    Direction consts.Direction
    Price     float64
    Volume    float64
    //status    consts.Status
    //datetime time.Time
}

func NewOrder(symbol string,
    exchange consts.Exchange,
    orderId int,
    direction consts.Direction,
    price float64,
    volume float64,
//status consts.Status,
//    datetime time.Time
) *Order {
    return &Order{
        Symbol:    symbol,
        Exchange:  exchange,
        OrderId:   orderId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        //status:    status,
        //datetime: datetime,
    }
}
