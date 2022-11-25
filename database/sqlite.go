package database

import (
	"fmt"
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	db *gorm.DB
}

var _sqlite *Sqlite

func NewSqlite() *Sqlite {
	if _sqlite != nil {
		return _sqlite
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&model.Bar{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_sqlite = &Sqlite{
		db: db,
	}
	return _sqlite

}

func (s *Sqlite) LoadBarData(
	symbol string,
	exchange consts.Exchange,
	interval consts.Interval,
	start time.Time,
	end time.Time,
) []model.Bar {
	var bars []model.Bar
	s.db.Where("symbol = ? AND exchange = ? AND interval = ? AND datetime >= ? AND datetime <= ?",
		symbol, exchange, interval, start, end).Order("datetime").Find(&bars)
	return bars
}

func (s *Sqlite) SaveBarData(bars []model.Bar) bool {
	tx := s.db.CreateInBatches(bars, 100)
	if tx.Error != nil {
		fmt.Println(tx.Error)
	}
	return tx.Error == nil
}
