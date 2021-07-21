package agin

var G = &config{}

func Init(configFile []byte) {
	G.ENV.Init()
	GetConfigFromYAML(configFile, G)
	G.init()
}

var _register = newAdminRegister()

var qiNiuUpload QiNiu

func SetQiNiu(qn QiNiu)  {
	qiNiuUpload = qn
}

func AdminRegisterModel(model IAdminModel, name string) {
	_register.AddAdminModel(model, name)
}