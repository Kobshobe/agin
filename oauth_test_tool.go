package agin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func OAuthTest(engine *gin.Engine, t *testing.T, qrUrl string, allowUrl string, token string, mode string) {
	var done = make(chan bool)
	var getQrOk = make(chan bool)
	var loginOk = false
	var testUuid = "qq"
	var head = http.Header{}

	if mode == "admin" {
		testUuid = "ww"
		head.Add("m", "admin")
	}

	s := httptest.NewServer(engine)
	//defer s.Close()

	wsUrl := fmt.Sprintf("ws%s%s", strings.TrimPrefix(s.URL, "http"), qrUrl)
	fmt.Println(wsUrl)
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, head)
	if err != nil {
		t.Fatalf("[ws test] create ws err: %v", err)
	}


	go func() {
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("[ws test] read msg err: %v", err)
		}
		_ = string(p)
		// 将二维码写入文件
		//file, err := os.OpenFile(fmt.Sprintf("./%s.jpeg", mode), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		//if err != nil {
		//	fmt.Println("open image err", err)
		//}
		//
		//_, err = file.Write(p)
		//if err != nil {
		//	t.Error("save qr image err")
		//}
		fmt.Println("[ws test] get QR img ok")
		getQrOk <- true
		//fmt.Println(response)
	}()

	go func() {
		<- getQrOk
		go func() {
			_, p, err := ws.ReadMessage()
			if err != nil {
				t.Fatalf("[ws test] read msg err: %v", err)
			}
			pStr := string(p)
			if strings.HasPrefix(pStr, `{"token":`) {
				fmt.Println("[ws test] get qToken info: ", pStr)
				loginOk = true
				tokenInfo := make(map[string]interface{})
				d := json.NewDecoder(bytes.NewReader(p))
				_ = d.Decode(&tokenInfo)
				//fmt.Println(tokenInfo)
				require.Equal(t, tokenInfo["liveTime"], G.WxApp.GetJwtLife().Hours())
				require.Equal(t, len(tokenInfo["token"].(string)) > 20, true)
				require.Equal(t, len(tokenInfo["openid"].(string)) > 10, true)
				done <- true
			} else {
				//fmt.Println(pStr)
				done <- false
			}
		}()

		// ==================== 模拟：扫描小程序码
		client := &http.Client{}
		request, err := http.NewRequest("GET", fmt.Sprintf("%s%s?uuid=%s",s.URL, allowUrl, testUuid), nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		request.Header.Add("Authorization", token)
		if mode == "admin" {
			request.Header.Add("m", "admin")
		}
		_, err = client.Do(request)
		if err != nil {
			t.Fatalf("[ws] allow login err: %v", err)
		}
	}()
	//defer func() {_ = ws.Close()}()
	if err = ws.WriteMessage(websocket.BinaryMessage, []byte("login")); err != nil {
		t.Fatalf("[ws test] write msg err: %v", err)
	}

	isOk := <-done
	if !isOk {
		t.Fatal("[ws test] get qToken err")
	}
	fmt.Println("get token ok")

	require.Equal(t, loginOk, true)
}