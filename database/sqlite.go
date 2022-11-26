package database

import (
	"fmt"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/model"
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
	db.AutoMigrate(&model.Bar{})
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
	symbol string,
	exchange consts.Exchange,
	interval consts.Interval,
	start string,
	end string,
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
