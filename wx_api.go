package agin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 客户端传来code获取openid和sessionKey
func WxGetOpenidAndSessionKey(code string, wxApp WxApp) (wxLoginInfo JWxLoginInfo, err error) {
	var (
		_url string
		resq *http.Response
		body []byte
	)
	wxLoginInfo = JWxLoginInfo{}
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
	//fmt.Println(string(body))
	if err != nil {
		return wxLoginInfo, errors.New("get skb err")
	}

	if err = json.Unmarshal(body, &wxLoginInfo); err != nil {
		return wxLoginInfo, errors.New("get sk js err")
	}
	if wxLoginInfo.Openid == "" {
		return wxLoginInfo, errors.New("empty Openid")
	}
	_ = resq.Body.Close()
	return wxLoginInfo, nil
}

type WxApp struct {
	Mode              string
	AppId             string `mapstructure:"appId" json:"appId" yaml:"appId"`
	AppSecret         string `mapstructure:"APPSECRET" json:"APPSECRET" yaml:"APPSECRET"`
	JwtLive           time.Duration
	JwtSecret         string `mapstructure:"jwtSecret" json:"jwtSecret" yaml:"jwtSecret"`
	TokenVersion      string `mapstructure:"tokenVersion" json:"tokenVersion" yaml:"tokenVersion"`
	AdminTokenVersion string `mapstructure:"adminTokenVersion" json:"adminTokenVersion" yaml:"adminTokenVersion"`

	AccessToken     string
	AccessTokenTime time.Time

	OAuthConnPool ConnectionPool

	adminConnPool ConnectionPool
	AdminOpenid []string

	upgrade websocket.Upgrader
}

func (w *WxApp) AddOAuthConn(write http.ResponseWriter, request *http.Request, head http.Header, mode string) (conn *Connection, err error) {
	var wsConn *websocket.Conn
	wsConn, err = w.upgrade.Upgrade(write, request, nil)
	if err != nil {
		fmt.Println("add conn create: ", err)
		return
	}

	if mode == "admin" {
		conn, err = w.adminConnPool.AddConnection(wsConn, mode)
	} else {
		conn, err = w.OAuthConnPool.AddConnection(wsConn, mode)
	}
	return
}

func (w *WxApp) GetOAuthConn(uuid string, mode string) (conn *Connection, ok bool) {
	if mode == "admin" {
		conn, ok = w.adminConnPool.Pool[uuid]
	} else {
		conn, ok = w.OAuthConnPool.Pool[uuid]
	}
	return
}

func (w WxApp) Upgrade(write http.ResponseWriter, request *http.Request, head http.Header) (conn *websocket.Conn, err error) {
	conn, err = w.upgrade.Upgrade(write, request, nil)

	return
}

type JWxLoginInfo struct {
	SessionKey string `json:"session_key"`
	Openid     string `json:"Openid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Code       string `json:"code"`
}

type JWxQr struct {
	Errcode     int    `json:"errcode"`
	QRBuffer    []byte `json:"buffer"`
	ContentType string `json:"contentType"`
}

type JWxAccessToken struct {
	AccessToken string `json:"access_token"`
	Errcode     int    `json:"errcode"`
	ContentType string `json:"contentType"`
}

func (w *WxApp) getAccessToken() (token string, err error) {
	if time.Since(w.AccessTokenTime) < time.Second*6500 {
		return w.AccessToken, nil
	}
	var jsonBody JWxAccessToken
	req, err := http.Get(w.GetAccessTokenUrl())
	if err != nil || req.StatusCode != 200 {
		return "1t err", err
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "2r err", err
	}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return "3j err", err
	}
	if jsonBody.Errcode != 0 {
		return "4g err", errors.New(fmt.Sprintf("errcode id = %d", jsonBody.Errcode))
	}
	//fmt.Println(w.AccessToken)
	w.AccessToken = jsonBody.AccessToken
	w.AccessTokenTime = time.Now()
	return "", nil
}

func (w WxApp) GetAccessTokenUrl() string {
	return fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		w.AppId, w.AppSecret)
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

func (w *WxApp) GetQRFromWX(scene string) (buffer []byte, error error) {
	var jsonBody JWxQr
	url, err := w.GetQRUrl()
	if err != nil {
		return nil, errors.New("qr get access token err: " + err.Error())
	}
	b := strings.NewReader(fmt.Sprintf(`{"scene":"%s","page":"pages/login/index"}`, scene))
	req, err := http.Post(url, "", b)
	defer req.Body.Close()
	if err != nil || req.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(body, &jsonBody) // 不要处理错误

	if jsonBody.Errcode != 0 {
		return nil, errors.New("get qr buffer err")
	}
	//fmt.Println(*w)
	//file, err := os.OpenFile("./222.jpeg", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	//if err != nil {
	//	fmt.Println("open image err", err)
	//	panic("open err")
	//}
	//defer file.Close()
	//file.Write(body)
	//fmt.Println(w.ContentType)
	return body, nil
}

func (w WxApp) GetQRUrl() (token string, err error) {
	token, err = w.getAccessToken()
	if err != nil {
		return
	}
	token = fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s",
		w.AccessToken)
	return
}
