package agin

import (
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type config struct {
	WxApp     *WxApp     `yaml:"wxApp"`
	Mysql     *Mysql     `yaml:"mysql"`
	System    *System    `yaml:"system"`
	Email     *Email     `yaml:"email"`
	Log       *Logger    `yaml:"logger"`
	GAnalytic *GAnalytic `yaml:"GAnalytic"`
	DB        *gorm.DB
	ENV       env
}

func (c *config) init() {
	if c.Mysql != nil {
		c.DB = c.Mysql.InitDB(G.System.Mode)
	}
	if c.WxApp != nil {
		c.WxApp.Init()
	}
	if c.Log != nil {
		c.Log.Init()
	}
	if c.Email != nil {
		c.Email.Init()
	}
	if c.GAnalytic != nil {
		c.GAnalytic.Init()
	}
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
