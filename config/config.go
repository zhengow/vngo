package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type config struct {
	*MysqlConfig   `yaml:"mysql"`
	*BinanceConfig `yaml:"binance"`
}

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

type BinanceConfig struct {
	Apikey    string `yaml:"apikey"`
	SecretKey string `yaml:"secretkey"`
}

func NewConfig(content []byte) (*config, error) {
	_config := &config{}
	err := yaml.Unmarshal(content, _config)
	if err != nil {
		return nil, err
	}
	return _config, nil
}
