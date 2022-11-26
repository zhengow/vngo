package database

import (
	"time"

	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/model"
)

type Database interface {
	LoadBarData(
		symbol string,
		exchange consts.Exchange,
		interval consts.Interval,
		start time.Time,
		end time.Time,
	) []model.Bar
	SaveBarData([]model.Bar) bool
}

var _db Database // default sqlite

func LoadBarData(
	symbol string,
	exchange consts.Exchange,
	interval consts.Interval,
	start time.Time,
	end time.Time,
) []model.Bar {
	return _db.LoadBarData(symbol, exchange, interval, start, end)
}
