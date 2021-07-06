package agin

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

type env struct {
	systemOS      string
	testStartTime time.Time
	RootPath      string
	StaticRoot    string
}

func (e *env) Init() {
	// go:embed

	if e.SystemOS() == "darwin" || G.ENV.SystemOS() == "windows" {
		_, p, _, _ := runtime.Caller(0)
		e.RootPath = p[:len(p)-13]
		log.Printf("-----working on development | os_env: %s | RootPath: %s ---------\n", e.SystemOS(), e.RootPath)
	} else {
		ePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		if os.Getenv("exercise_test") == "1" {
			_, p, _, _ := runtime.Caller(0)
			e.RootPath = p[:len(p)-13]
		} else {
			e.RootPath = path.Dir(ePath)
		}
		gin.SetMode(gin.ReleaseMode)
		//log.Printf("-----working on production | os_env: %s | dns_name:%s---------\n", e.SystemOS())
	}
	e.StaticRoot = e.RootPath + "/static"
}

func (e *env) SetSystemOS() {
	e.systemOS = runtime.GOOS
}

func (e env) SystemOS() string {
	return runtime.GOOS
}