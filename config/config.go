package config

import "gopkg.in/yaml.v2"

type config struct {
	*MysqlConfig `yaml:"mysql"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func NewConfig(content []byte) (*config, error) {
	_config := &config{}
	err := yaml.Unmarshal(content, _config)
	if err != nil {
		return nil, err
	}
	return _config, nil
}
