package agin

import (
	"time"
)

var G = &config{}

func Init(configFile []byte) {
	G.ENV.Init()
	GetConfigFromYAML(configFile, G)
	if G.Mysql != nil {
		G.DB = G.Mysql.InitDB(G.System.Mode)
	}
	if G.WxApp != nil {
		G.WxApp.JwtLive = time.Hour * 24 * 14
		G.WxApp.Init()
	}
}

var _register = newAdminRegister()

var qiNiuUpload QiNiu

func SetQiNiu(qn QiNiu)  {
	qiNiuUpload = qn
}

func AdminRegisterModel(model IAdminModel, name string) {
	_register.AddAdminModel(model, name)
}