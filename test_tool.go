package agin

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewTHttp(router *gin.Engine,t *testing.T, url string, method string, body io.Reader) THttp {
	req, _ := http.NewRequest(method, url, body)
	return THttp{Router: router, T: t, Url: url, Method: method, RequestBody: body, Writer: httptest.NewRecorder(), Request: req}
}

type THttp struct {
	Router  		*gin.Engine
	T       		*testing.T
	Url     		string
	Method  		string
	RequestBody    	io.Reader
	Writer  		*httptest.ResponseRecorder
	Request 		*http.Request
}


func (h *THttp) BaseTest(status int, bodyString string, message ...interface{}) {
	h.Router.ServeHTTP(h.Writer, h.Request)

	m := []interface{}{h.Method + " url: " + h.Url + "\n"}
	m = append(m, message, "ResponseBody:"+h.Writer.Body.String())

	assert.Equal(h.T, status, h.Writer.Code, m...)

	if bodyString != "" {
		assert.Equal(h.T, bodyString, h.Writer.Body.String(), m...)
	}
}


func (h *THttp) GetList() (result []interface{}) {
	d := json.NewDecoder(h.Writer.Body)
	d.UseNumber()
	_ = d.Decode(&result)
	return
}

func (h *THttp) GetMap() (result map[string]interface{}) {
	d := json.NewDecoder(h.Writer.Body)
	d.UseNumber()
	_ = d.Decode(&result)
	return
}
