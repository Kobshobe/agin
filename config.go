package agin

import (
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"time"
)

type config struct {
	WxApp  *WxApp  `yaml:"wxApp"`
	Mysql  *Mysql  `yaml:"mysql"`
	System *System `yaml:"system"`
	Email  *Email  `yaml:"email"`
	Log    *Logger  `yaml:"logger"`
	DB     *gorm.DB
	ENV    env
}

func (c *config) Init() {
	if c.Mysql != nil {
		c.DB = c.Mysql.InitDB(G.System.Mode)
	}
	if c.WxApp != nil {
		c.WxApp.JwtLive = time.Hour * 24 * 14
		c.WxApp.Init()
	}
	if c.Log != nil {
		c.Log.Init()
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
