package agin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type Event struct {
	Client        string
	ClientVersion string
	Name          string
	IP            string
}

type GAnalytic struct {
	MeasurementId string `yaml:"measurementId"`
	ApiSecret     string `yaml:"apiSecret"`
	ClientID      string `yaml:"clientId"`
	SystemInfoPer int 	`yaml:"systemInfoPer"`
	url           string
}

func (ga *GAnalytic) Init() {
	ga.url = fmt.Sprintf("https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s", ga.MeasurementId, ga.ApiSecret)
	ga.SendSystemInfo()
}

func (ga *GAnalytic) Send(name string) {

}

func (ga *GAnalytic) SendSystemInfo() {
	ticker := time.NewTicker(time.Minute * 20)

	go func() {
		for {
			select {
			case <-ticker.C:
				s := ga.getSystemInfo()
				//fmt.Println(s)
				info := &AnalyticInfo{
					ClientID: ga.ClientID,
					Events: AnalyticEvents{
						Name: "server_system_info",
						Params: s,
					},
				}
				body, err := json.Marshal(info)
				if err != nil {
					fmt.Println(err)
					if G.Email != nil {
						G.Email.SendToDefault("json marshal system info err", "")
					}
					return
				}
				_, _ = http.Post(ga.url, "application/json", bytes.NewReader(body))
			}
		}
	}()


}

func (ga GAnalytic) getSystemInfo() map[string]interface{} {
	start := time.Now()
	MemStats := &runtime.MemStats{}
	runtime.ReadMemStats(MemStats)
	return map[string]interface{}{
		"NumGoroutine": runtime.NumGoroutine(),  //
		"Alloc": MemStats.Alloc/1000000, // golang语言框架堆空间分配的字节数
		"TotalAlloc": MemStats.TotalAlloc/1000000,  // 从服务开始运行至今分配器为分配的堆空间总
		"Sys": MemStats.Sys/1000000, //服务现在系统使用的内存
		"OtherSys": MemStats.OtherSys/1000000,  //golang系统架构占用的额外空间
		"Lookups": MemStats.Lookups,  // 被runtime监视的指针数
		"PauseTotalNs": MemStats.PauseTotalNs/1000000,  //
		"NumGC": MemStats.NumGC,  // 垃圾回收调用次数
		"cost": time.Since(start),  //
	}
}

type AnalyticEvents struct {
	Name   string	`json:"name"`
	Params interface{} `json:"params"`
}

type SystemInfo struct {
	NumGoroutine int
}

type AnalyticInfo struct {
	ClientID string	`json:"client_id"`
	Events   AnalyticEvents	`json:"events"`
}
