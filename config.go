package agin

import (
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type config struct {
	WxApp  *WxApp `yaml:"wxApp"`
	Mysql  *Mysql `yaml:"mysql"`
	System *System `yaml:"system"`
	Email  *Email `yaml:"email"`
	DB *gorm.DB
	ENV env
}

type System struct {
	Mode string `yaml:"mode"`
}

func GetConfigFromYAML(configSource []byte, config interface{}) {
	err := yaml.Unmarshal(configSource, config)
	if err != nil {
		panic(err)
	}
}