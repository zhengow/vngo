package consts

type Interval string
type Exchange string

type interval struct {
	MINUTE Interval
    HOUR Interval
    DAILY Interval
    WEEKLY Interval
    TICK Interval
    TRANSACTION Interval
}

var IntervalEnum = interval {
	MINUTE: "1m",
    HOUR: "1h",
    DAILY: "d",
    WEEKLY: "w",
}

type exchange struct {
    BINANCE Exchange
}

var ExchangeEnum = exchange {
    BINANCE: "BINANCE",
}

const DateFormat = "2006-01-02 15:04:05"