package vngo

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Mysql struct {
    db *gorm.DB
}

var _mysql *Mysql

func NewMysql(mysqlConfig *config.MysqlConfig) {
    if _mysql != nil {
        _db = _mysql
    }
    // user := mysqlConfig.User
    // password := mysqlConfig.Password
    // port := mysqlConfig.Port
    // host := mysqlConfig.Host
    // dbName := mysqlConfig.Name
    // dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
    dsn := mysqlConfig.GetDsn()
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
    symbol Symbol,
    interval Interval,
    start string,
    end string,
) []Bar {
    var bars []Bar
    s.db.Table("dbbardata").Where(fmt.Sprintf("symbol='%s' AND exchange='%s' AND `interval`='%s' AND datetime>='%s' AND datetime<='%s'",
        symbol.Symbol, symbol.Exchange, interval, start, end)).Order("datetime").Find(&bars)
    return bars
}

func (s *Mysql) SaveBarData(bars []Bar) bool {
    tx := s.db.CreateInBatches(bars, 100)
    if tx.Error != nil {
        fmt.Println(tx.Error)
    }
    return tx.Error == nil
}
