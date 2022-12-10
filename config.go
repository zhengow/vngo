package vngo

import (
    "github.com/zhengow/vngo/config"
    "gopkg.in/yaml.v2"
)

type _config struct {
    *config.MysqlConfig   `yaml:"mysql"`
    *config.BinanceConfig `yaml:"binance"`
}

func NewConfig(content []byte) (*_config, error) {
    cfg := &_config{}
    err := yaml.Unmarshal(content, cfg)
    if err != nil {
        return nil, err
    }
    return cfg, nil
}
