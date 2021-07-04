package agin

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var _register = NewRegister()

var adminDB *gorm.DB

var wxApp *WxApp

var qiNiuUpload QiNiu

func SetDB(_db *gorm.DB)  {
	adminDB = _db
}

func SetQiNiu(qn QiNiu)  {
	qiNiuUpload = qn
}

func AdminRegisterModel(model IAdminModel, name string) {
	_register.AddAdminModel(model, name)
}

func SetWxApp(_wxApp *WxApp) {
	wxApp = _wxApp
	wxApp.upgrade = websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		//ReadBufferSize:  1,
		//WriteBufferSize: 1,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}