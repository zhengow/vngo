package config

import "fmt"

type MysqlConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DbName   string `yaml:"dbname"`
}

func (c *MysqlConfig) GetDsn() string {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DbName)
    return dsn
}
