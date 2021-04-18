package agin_router

import "github.com/gin-gonic/gin"


// 路由注册器
type engine struct {
	router *gin.Engine
}

func (e engine) New() engine {
	return engine{router: gin.Default()}
}

func (e engine) GetRouter() *gin.Engine {
	return e.router
}

func (e *engine) AddGroup(register func(engine *gin.Engine, urlPrefix string), urlPrefix string) {
	register(e.router, urlPrefix)
}
