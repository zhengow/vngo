package consts

type Interval string
type Exchange string
type Direction string
type Status string

type interval struct {
    MINUTE      Interval
    HOUR        Interval
    DAILY       Interval
    WEEKLY      Interval
    TICK        Interval
    TRANSACTION Interval
}

var IntervalEnum = interval{
    MINUTE: "1m",
    HOUR:   "1h",
    DAILY:  "d",
    WEEKLY: "w",
}

type exchange struct {
    BINANCE Exchange
}

var ExchangeEnum = exchange{
    BINANCE: "BINANCE",
}

type direction struct {
    LONG  Direction
    SHORT Direction
}

var DirectionEnum = direction{
    LONG:  "LONG",
    SHORT: "SHORT",
}

type status struct {
    SUBMITTING Status
    NOTTRADED  Status
    PARTTRADED Status
    ALLTRADED  Status
    CANCELLED  Status
    REJECTED   Status
}

var StatusEnum = status{
    SUBMITTING: "提交中",
    NOTTRADED:  "未成交",
    PARTTRADED: "部分成交",
    ALLTRADED:  "全部成交",
    CANCELLED:  "已撤销",
    REJECTED:   "拒单",
}

const DateFormat = "2006-01-02 15:04:05"
