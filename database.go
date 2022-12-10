package vngo

import (
    "fmt"
)

type Database interface {
    LoadBarData(
        symbol Symbol,
        interval Interval,
        start string,
        end string,
    ) []Bar
    SaveBarData([]Bar) bool
}

var _db Database // default sqlite

func LoadBarData(
    symbol Symbol,
    interval Interval,
    start string,
    end string,
) []Bar {
    if _db == nil {
        fmt.Println("db is nil")
        return nil
    }
    return _db.LoadBarData(symbol, interval, start, end)
}
