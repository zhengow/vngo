package strategy

import (
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/types"
)

type Order struct {
    Symbol    Symbol
    OrderId   string
    Direction types.Direction
    Price     float64
    Volume    float64
    Status    types.Status
    Traded    float64
    //datetime strategy.VnTime
}

func NewOrder(symbol Symbol,
    orderId string,
    direction types.Direction,
    price float64,
    volume float64,
//status Status,
//    datetime strategy.VnTime
) *Order {
    return &Order{
        Symbol:    symbol,
        OrderId:   orderId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        Status:    consts.StatusEnum.SUBMITTING,
        //datetime: datetime,
    }
}
