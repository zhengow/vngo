package database

import (
	"fmt"

	"github.com/zhengow/vngo/config"
	"github.com/zhengow/vngo/consts"
	"github.com/zhengow/vngo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	db *gorm.DB
}

var _mysql *Mysql

func NewMysql() {
	if _mysql != nil {
		_db = _mysql
	}
	mysqlConfig := config.Config.Mysql
	user := mysqlConfig.User
	password := mysqlConfig.Password
	port := mysqlConfig.Port
	host := mysqlConfig.Host
	dbName := mysqlConfig.Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	_mysql = &Mysql{
		db: db,
	}
	_db = _mysql
}

func (s *Mysql) LoadBarData(
	symbol string,
	exchange consts.Exchange,
	interval consts.Interval,
	start string,
	end string,
) []model.Bar {
	var bars []model.Bar
	s.db.Table("dbbardata").Where(fmt.Sprintf("symbol='%s' AND exchange='%s' AND `interval`='%s' AND datetime>='%s' AND datetime<='%s'", 
			symbol, exchange, interval, start, end)).Order("datetime").Find(&bars)
	return bars
}

func (s *Mysql) SaveBarData(bars []model.Bar) bool {
	tx := s.db.CreateInBatches(bars, 100)
	if tx.Error != nil {
		fmt.Println(tx.Error)
	}
	return tx.Error == nil
}
