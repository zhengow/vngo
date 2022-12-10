package vngo

import (
    "github.com/zhengow/vngo/config"
    "github.com/zhengow/vngo/database"
)

func UseMysql(mysqlConfig *config.MysqlConfig) {
    database.DB = database.NewMysql(mysqlConfig)
}
