package database

import (
    "fmt"
    "github.com/zhengow/vngo"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Mysql struct {
    db *gorm.DB
}

func NewMysql(mysqlConfig *vngo.MysqlConfig) *Mysql {
    dsn := mysqlConfig.GetDsn()
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Println(err)
        return nil
    }
    return &Mysql{
        db: db,
    }
}

func (s *Mysql) LoadBarData(
    symbol vngo.Symbol,
    interval vngo.Interval,
    start string,
    end string,
) []vngo.Bar {
    var bars []vngo.Bar
    s.db.Table("dbbardata").Where(fmt.Sprintf("symbol='%s' AND exchange='%s' AND `interval`='%s' AND datetime>='%s' AND datetime<='%s'",
        symbol.Name, symbol.Exchange, interval, start, end)).Order("datetime").Find(&bars)
    return bars
}

func (s *Mysql) SaveBarData(bars []vngo.Bar) bool {
    tx := s.db.CreateInBatches(bars, 100)
    if tx.Error != nil {
        fmt.Println(tx.Error)
    }
    return tx.Error == nil
}
