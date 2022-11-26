package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

//go:embed dev.yml
var content []byte

type config struct {
	Mysql *Mysql `yaml:"mysql"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

var Config *config

func init() {
	err := yaml.Unmarshal(content, &Config)
	if err != nil {
		fmt.Println(err)
	}
}
