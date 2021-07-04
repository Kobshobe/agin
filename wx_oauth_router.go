package agin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)


func WxOAuthRouterRegister(router *gin.RouterGroup, qrUrl string, allowUrl string) {
	if wxApp.JwtSecret == "" || wxApp.AppId == "" || wxApp.Mode == "" {
		panic("Should init wxApp")
	}
	wxApp.OAuthConnPool = NewConnectionPool()
	router.GET(qrUrl, QrLogin)
	router.GET(allowUrl, allowLogin)
}


func WxAdminOAuthRouterRegister(router *gin.RouterGroup, qrUrl string, allowUrl string) {
	if wxApp.JwtSecret == "" || wxApp.AppId == "" || wxApp.Mode == "" {
		panic("Should init wxApp")
	}
	wxApp.adminConnPool = NewConnectionPool()
	router.GET(qrUrl, QrLogin)
	router.GET(allowUrl, allowLogin)
}


func QrLogin(c *gin.Context) {
	fmt.Println("QrLogin", c.FullPath())
	var (
		err    error
		conn   *Connection
	)

	if conn, err = wxApp.AddOAuthConn(c.Writer, c.Request, nil, c.GetHeader("m")); err != nil {
		goto ERR
	}

	for {
		if _, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		//if err = conn.WriteMessage(data); err != nil {
		//	goto ERR
		//}
	}

ERR:
	if conn != nil {
		conn.Close()
	}

}


func allowLogin(c *gin.Context) {
	var (
		uuid string
		conn *Connection
		ok   bool
		err error
		openid string
	)

	uuid = c.Query("uuid")
	conn, ok = wxApp.GetOAuthConn(uuid, c.GetHeader("m"))
	if ok == false || conn.isClosed == true {
		fmt.Println("no in pool")
		c.JSON(200, gin.H{
			"err": "ok",
			"msg": "invalid qr",
		})
		return
	}

	// 生成token
	openid, _, err = wxApp.GetTokenInfo(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{
			"err": "err Token",
		})
		_ = c.Error(errors.New("err Token"))
		return
	}
	err = conn.SetToken(openid)
	if err != nil {
		c.JSON(200, gin.H{
			"err": "ok",
			"msg": "new t err",
		})
		return
	}
	conn.Openid = openid
	//fmt.Println(Token)

	conn.allowLoginChan <- "allow"
	for {
		msg := <- conn.loginResultChan
		//fmt.Println("LoginResultChan: ", msg)
		if msg == "loginOk" {
			//fmt.Println("ok login ok")
			c.JSON(200, gin.H{
				"err": "ok",
				"msg": "loginOk",
			})
			break
		}
	}
	//c.JSON(200, gin.H{
	//	"err": "ok",
	//	"msg": "qr login ok",
	//})
}
