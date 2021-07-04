package agin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRouterRegister(r *gin.Engine, urlPrefix string, db *gorm.DB) {
	SetDB(db)
	group := r.Group(urlPrefix)
	group.Use(MAdminAuth())
	group.Use(Cors())
	{
		group.GET("/model", getModelsList)
		group.POST("/get_record/:table", getRecords)
		group.PUT("/record/:table", updateRecord)
		group.POST("/record/:table", createRecord)
		group.DELETE("/record/:table", deleteRecord)
		group.POST("/upload_file", uploadFile)
	}
	loginGroup := r.Group(urlPrefix)
	WxAdminOAuthRouterRegister(loginGroup, "/login", "/allow_login")
}