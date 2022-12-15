package models

import (
    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/types"
)

type TradeData struct {
    Symbol    Symbol
    OrderId   string
    TradeId   int
    Direction types.Direction
    Price     float64
    Volume    float64
    Datetime  VnTime
}

func NewTradeData(symbol Symbol,
    orderId string,
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
    return t.Direction == consts.DirectionEnum.SHORT
}
