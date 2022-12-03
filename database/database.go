package database

import (
    "fmt"

    "github.com/zhengow/vngo/consts"
    "github.com/zhengow/vngo/model"
)

type Database interface {
    LoadBarData(
        symbol model.Symbol,
        interval consts.Interval,
        start string,
        end string,
    ) []model.Bar
    SaveBarData([]model.Bar) bool
}

var _db Database // default sqlite

func LoadBarData(
    symbol model.Symbol,
    interval consts.Interval,
    start string,
    end string,
) []model.Bar {
    if _db == nil {
        fmt.Println("db is nil")
        return nil
    }
    return _db.LoadBarData(symbol, interval, start, end)
}
