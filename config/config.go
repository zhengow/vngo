package config

import "gopkg.in/yaml.v2"

type Config struct {
	*MysqlConfig `yaml:"mysql"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

func NewConfig(content []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
