package vngo

import (
    "fmt"
    "github.com/zhengow/vngo/config"
    "github.com/zhengow/vngo/database"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/types"
)

type Database interface {
    LoadBarData(
        symbol models.Symbol,
        interval types.Interval,
        start string,
        end string,
    ) []models.Bar
    SaveBarData([]models.Bar) bool
}

var _db Database // default sqlite

func init() {
    _db = database.NewSqlite()
}

func LoadBarData(
    symbol models.Symbol,
    interval types.Interval,
    start string,
    end string,
) []models.Bar {
    if _db == nil {
        fmt.Println("db is nil")
        return nil
    }
    return _db.LoadBarData(symbol, interval, start, end)
}

func UseMysql(mysqlConfig *config.MysqlConfig) {
    _db = database.NewMysql(mysqlConfig)
}
