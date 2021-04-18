package agin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)


//func Login(c *gin.Context) {
//	code := c.Query("code")
//	if code == "" {
//		c.JSON(400, gin.H{
//			"err": "empty code",
//		})
//		return
//	}
//	authInfo, err := GetOpenidAndSessionKey(code)
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(400, gin.H{
//			"err": "invalid code",
//		})
//		return
//	}
//	_user.Openid = authInfo.Openid
//
//	result := global.DB.FirstOrInit(&User{}, &User{Openid: authInfo.Openid})
//	if result.RowsAffected == 0 {
//		result = global.DB.Create(&_user)
//		if result.Error != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"err":"create user err",
//			})
//			return
//		}
//	}
//	sk := SessionAuth{SessionKey: authInfo.SessionKey, Openid: authInfo.Openid}
//	global.DB.Create(&sk)
//	if result.Error != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"err":"create user sk err",
//		})
//		return
//	}
//	c.JSON(200, gin.H{
//		"err": "ok",
//		"sessionKey": authInfo.SessionKey,
//		"openid": authInfo.Openid,
//	})
//}

// 客户端传来code获取openid和sessionKey
func WxGetOpenidAndSessionKey(code string, wxApp WxApp) (wxLoginInfo WxLoginInfo, err error) {
	var (
		_url  string
		resq  *http.Response
		body []byte
	)
	wxLoginInfo = WxLoginInfo{}
	if code == "" {
		err = errors.New("empty code")
		return
	}
	_url = wxApp.WxLoginUrl(code)
	resq, err = http.Get(_url)

	if err != nil {
		return wxLoginInfo, errors.New("get sk err")
	}
	body, err = ioutil.ReadAll(resq.Body)
	fmt.Println(string(body))
	if err != nil {
		return wxLoginInfo, errors.New("get skb err")
	}

	if err = json.Unmarshal(body, &wxLoginInfo); err != nil {
		return wxLoginInfo, errors.New("get sk js err")
	}
	if wxLoginInfo.Openid == "" {
		return wxLoginInfo, errors.New("empty openid")
	}
	_ = resq.Body.Close()
	return wxLoginInfo, nil
}

type WxApp struct {
	AppId 			string			`mapstructure:"appId" json:"appId" yaml:"appId"`
	AppSecret 		string			`mapstructure:"APPSECRET" json:"APPSECRET" yaml:"APPSECRET"`
	JwtLive			time.Duration
	JwtSecret 		string			`mapstructure:"jwtSecret" json:"jwtSecret" yarm:"jwtSecret"`
}

type WxLoginInfo struct {
	SessionKey string	`json:"session_key"`
	Openid string		`json:"openid"`
	ErrCode int			`json:"errcode"`
	ErrMsg string		`json:"errmsg"`
	Code string			`json:"code"`
}

func (w WxLoginInfo) NewJwtToken(versionInfo string, wxApp WxApp) (token string, err error) {
	if w.Openid == "" {
		err = errors.New("err openid")
		return
	}
	token, err = CreateJwtToken(
		getJwtUid(w.Openid, versionInfo),
		wxApp.JwtSecret, wxApp.JwtLive)
	if err != nil {
		err = errors.New("c token err")
	}
	return
}

func (w WxApp) GetTokenInfo(token string) (openid string, versionInfo string, err error) {
	info, err := ParseJwtToken(token, w.JwtSecret)
	if err != nil {
		return
	}
	l := strings.Split(info, "##")
	openid = l[0]
	versionInfo = l[1]
	return
}

func (w WxApp) WxLoginUrl(code string) string {
	return fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		w.AppId, w.AppSecret, code)
}

func (w WxApp) NewJwtToken(openid, versionInfo string) (token string, err error) {
	token, err = CreateJwtToken(getJwtUid(openid, versionInfo), w.JwtSecret, w.JwtLive)

	return
}

func getJwtUid(openid string, versionInfo string) (uid string) {
	uid = fmt.Sprintf("%s##%s", openid, versionInfo)
	return
}