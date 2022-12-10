package vngo

import (
    "fmt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Sqlite struct {
    db *gorm.DB
}

var _sqlite *Sqlite

func NewSqlite() {
    if _sqlite != nil {
        _db = _sqlite
    }
    db, err := gorm.Open(sqlite.Open("vngo.db"), &gorm.Config{})
    db.AutoMigrate(&Bar{})
    if err != nil {
        fmt.Println(err)
        return
    }
    _sqlite = &Sqlite{
        db: db,
    }
    _db = _sqlite
}

func (s *Sqlite) LoadBarData(
    symbol Symbol,
    interval Interval,
    start string,
    end string,
) []Bar {
    var bars []Bar
    s.db.Where("symbol = ? AND exchange = ? AND interval = ? AND datetime >= ? AND datetime <= ?",
        symbol.Symbol, symbol.Exchange, interval, start, end).Order("datetime").Find(&bars)
    return bars
}

func (s *Sqlite) SaveBarData(bars []Bar) bool {
    tx := s.db.CreateInBatches(bars, 100)
    if tx.Error != nil {
        fmt.Println(tx.Error)
    }
    return tx.Error == nil
}
