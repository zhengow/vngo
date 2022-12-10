package consts

import "github.com/zhengow/vngo/types"

type interval struct {
    MINUTE      types.Interval
    HOUR        types.Interval
    DAILY       types.Interval
    WEEKLY      types.Interval
    TICK        types.Interval
    TRANSACTION types.Interval
}

var IntervalEnum = interval{
    MINUTE: "1m",
    HOUR:   "1h",
    DAILY:  "d",
    WEEKLY: "w",
}

type exchange struct {
    BINANCE types.Exchange
}

var ExchangeEnum = exchange{
    BINANCE: "BINANCE",
}

type direction struct {
    LONG  types.Direction
    SHORT types.Direction
}

var DirectionEnum = direction{
    LONG:  "LONG",
    SHORT: "SHORT",
}

type status struct {
    SUBMITTING types.Status
    NOTTRADED  types.Status
    PARTTRADED types.Status
    ALLTRADED  types.Status
    CANCELLED  types.Status
    REJECTED   types.Status
}

var StatusEnum = status{
    SUBMITTING: "提交中",
    NOTTRADED:  "未成交",
    PARTTRADED: "部分成交",
    ALLTRADED:  "全部成交",
    CANCELLED:  "已撤销",
    REJECTED:   "拒单",
}
