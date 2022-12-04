package model

import (
    "github.com/zhengow/vngo/consts"
    "time"
)

type TradeData struct {
    Symbol    Symbol
    OrderId   int
    TradeId   int
    Direction consts.Direction
    Price     float64
    Volume    float64
    Datetime  time.Time
}

func NewTradeData(symbol Symbol,
    orderId int,
    tradeId int,
    direction consts.Direction,
    price float64,
    volume float64,
    datetime time.Time,
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

func (t *TradeData) IsBuy() bool {
    return t.Direction == consts.DirectionEnum.LONG
}

func (t *TradeData) IsSell() bool {
    return t.Direction == consts.DirectionEnum.SHORT
}
