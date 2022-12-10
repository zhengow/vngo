package database

import (
    "fmt"
    "github.com/zhengow/vngo/strategy"
    "github.com/zhengow/vngo/types"
)

type Database interface {
    LoadBarData(
        symbol strategy.Symbol,
        interval types.Interval,
        start string,
        end string,
    ) []strategy.Bar
    SaveBarData([]strategy.Bar) bool
}

var DB Database // default sqlite

func LoadBarData(
    symbol strategy.Symbol,
    interval types.Interval,
    start string,
    end string,
) []strategy.Bar {
    if DB == nil {
        DB = NewSqlite()
        if DB == nil {
            fmt.Println("init sqlite error")
            return nil
        }
    }
    return DB.LoadBarData(symbol, interval, start, end)
}
