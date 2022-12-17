package database

import (
    "fmt"
    "github.com/zhengow/vngo/models"
)

type Database interface {
    LoadBarData(
        symbol models.Symbol,
        interval models.Interval,
        start string,
        end string,
    ) []models.Bar
    SaveBarData([]models.Bar) bool
}

var DB Database // default sqlite

func LoadBarData(
    symbol models.Symbol,
    interval models.Interval,
    start string,
    end string,
) []models.Bar {
    if DB == nil {
        DB = NewSqlite()
        if DB == nil {
            fmt.Println("init sqlite error")
            return nil
        }
    }
    return DB.LoadBarData(symbol, interval, start, end)
}
