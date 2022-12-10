package database

import (
    "fmt"
    "github.com/zhengow/vngo/config"
    "github.com/zhengow/vngo/models"
    "github.com/zhengow/vngo/types"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Mysql struct {
    db *gorm.DB
}

func NewMysql(mysqlConfig *config.MysqlConfig) *Mysql {
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
    symbol models.Symbol,
    interval types.Interval,
    start string,
    end string,
) []models.Bar {
    var bars []models.Bar
    s.db.Table("dbbardata").Where(fmt.Sprintf("symbol='%s' AND exchange='%s' AND `interval`='%s' AND datetime>='%s' AND datetime<='%s'",
        symbol.Name, symbol.Exchange, interval, start, end)).Order("datetime").Find(&bars)
    return bars
}

func (s *Mysql) SaveBarData(bars []models.Bar) bool {
    tx := s.db.CreateInBatches(bars, 100)
    if tx.Error != nil {
        fmt.Println(tx.Error)
    }
    return tx.Error == nil
}
