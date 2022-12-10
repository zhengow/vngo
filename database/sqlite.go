package database

import (
    "fmt"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/types"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Sqlite struct {
    db *gorm.DB
}

func NewSqlite() *Sqlite {
    db, err := gorm.Open(sqlite.Open("vngo.db"), &gorm.Config{})
    db.AutoMigrate(&models.Bar{})
    if err != nil {
        fmt.Println(err)
        return nil
    }
    return &Sqlite{
        db: db,
    }
}

func (s *Sqlite) LoadBarData(
    symbol models.Symbol,
    interval types.Interval,
    start string,
    end string,
) []models.Bar {
    var bars []models.Bar
    s.db.Where("symbol = ? AND exchange = ? AND interval = ? AND datetime >= ? AND datetime <= ?",
        symbol.Name, symbol.Exchange, interval, start, end).Order("datetime").Find(&bars)
    return bars
}

func (s *Sqlite) SaveBarData(bars []models.Bar) bool {
    tx := s.db.CreateInBatches(bars, 100)
    if tx.Error != nil {
        fmt.Println(tx.Error)
    }
    return tx.Error == nil
}
