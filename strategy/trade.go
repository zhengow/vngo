package strategy

import (
    "github.com/zhengow/vngo/enum"
    "github.com/zhengow/vngo/types"
)

type TradeData struct {
    Symbol    Symbol
    OrderId   int
    TradeId   int
    Direction types.Direction
    Price     float64
    Volume    float64
    Datetime  VnTime
}

func NewTradeData(symbol Symbol,
    orderId int,
    tradeId int,
    direction types.Direction,
    price float64,
    volume float64,
    datetime VnTime,
) *TradeData {
    return &TradeData{
        Symbol:    symbol,
        OrderId:   orderId,
        TradeId:   tradeId,
        Direction: direction,
        Price:     price,
        Volume:    volume,
        Datetime:  datetime,
    }
}

func (t *TradeData) IsSell() bool {
    return t.Direction == enum.DirectionEnum.SHORT
}
