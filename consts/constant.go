package consts

type interval struct {
	MINUTE string
    HOUR string
    DAILY string
    WEEKLY string
    TICK string
    TRANSACTION string
}

var Interval = interval {
	MINUTE: "1m",
    HOUR: "1h",
    DAILY: "d",
    WEEKLY: "w",
}