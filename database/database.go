package database

import (
    "fmt"
    "github.com/zhengow/vngo"
)

type Database interface {
    LoadBarData(
        symbol vngo.Symbol,
        interval vngo.Interval,
        start string,
        end string,
    ) []vngo.Bar
    SaveBarData([]vngo.Bar) bool
}

var _db Database // default sqlite

func LoadBarData(
    symbol vngo.Symbol,
    interval vngo.Interval,
    start string,
    end string,
) []vngo.Bar {
    if _db == nil {
        fmt.Println("db is nil")
        return nil
    }
    return _db.LoadBarData(symbol, interval, start, end)
}
