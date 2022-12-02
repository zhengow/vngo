package model

import (
    "github.com/zhengow/vngo/consts"
    "time"
)

type TradeData struct {
    symbol    string
    exchange  consts.Exchange
    orderId   int
    tradeId   int
    direction consts.Direction
    price     float64
    volume    float64
    datetime  time.Time
}

func NewTradeData(symbol string,
    exchange consts.Exchange,
    orderId int,
    tradeId int,
    direction consts.Direction,
    price float64,
    volume float64,
    datetime time.Time,
) *TradeData {
    return &TradeData{
        symbol:    symbol,
        exchange:  exchange,
        orderId:   orderId,
        tradeId:   tradeId,
        direction: direction,
        price:     price,
        volume:    volume,
        datetime:  datetime,
    }
}
