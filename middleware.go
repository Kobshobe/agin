package agin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/user"
	"time"
)

func MLogger(filePath string, isTest bool) gin.HandlerFunc {
	if isTest {
		u, _ := user.Current()
		filePath = u.HomeDir + filePath
	}

	linkName := "latest_log.log"
	scr, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic("open log file err:" + err.Error())
	}

	logger := logrus.New()

	logger.Out = scr

	logger.SetLevel(logrus.DebugLevel)

	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log",
		retalog.WithMaxAge(360*24*time.Hour),
		retalog.WithRotationTime(7*24*time.Hour),
		retalog.WithLinkName(linkName),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	Hook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})

	logger.AddHook(Hook)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		spendTime := time.Since(startTime).Microseconds()
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func MAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			openid string
			err    error
			token string
			isAllow = false
		)

		token = c.Request.Header.Get("Authorization")
		if token == "" {
			err = errors.New("t empty")
			goto ERR
		}

		openid, _, err = G.WxApp.GetTokenInfo(token)
	ERR:
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(401, gin.H{
				"err": err.Error(),
			})
			c.Abort()
			return
		}

		for _, id := range G.WxApp.AdminOpenid {
			if openid == id {
				isAllow = true
			}
		}

		if !isAllow {
			err = errors.New("no allow")
			goto ERR
		}

		c.Set("openid", openid)
	}
}

func MWxAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			openid string
			err    error
			token string
		)

		token = c.Request.Header.Get("Authorization")
		if token == "" {
			err = errors.New("t empty")
			goto ERR
		}

		openid, _, err = G.WxApp.GetTokenInfo(token)
	ERR:
		if err != nil {
			c.JSON(401, gin.H{
				"err": "t err",
			})
			c.Abort()
			return
		}

		c.Set("openid", openid)
	}
}
