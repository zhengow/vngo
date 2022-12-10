package models

import "github.com/zhengow/vngo/types"

type Order struct {
    Symbol    Symbol
    OrderId   int
    Direction types.Direction
    Price     float64
    Volume    float64
    //status    Status
    //datetime models.VnTime
}

func NewOrder(symbol Symbol,
    orderId int,
    direction types.Direction,
    price float64,
    volume float64,
//status Status,
//    datetime models.VnTime
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
