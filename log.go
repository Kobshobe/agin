package agin

import (
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/user"
	"time"
)

func logger(filePath string, isTest bool) gin.HandlerFunc {
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


type Logger struct {
	logger *zap.Logger
	Sugar *zap.SugaredLogger

	Filename string		`yaml:"filename"`
	MaxSize int
	MaxAge int
	MaxBackups int
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *Logger) Init() {
	core := zapcore.NewCore(getEncoder(), l.getLogWriter(), zap.DebugLevel)
	l.logger = zap.New(core)
	l.Sugar = l.logger.Sugar()
}

func (l Logger) getLogWriter() zapcore.WriteSyncer {
	if G.System.Mode == "test" {
		u, _ := user.Current()
		l.Filename = u.HomeDir + l.Filename
		_, err := os.OpenFile(l.Filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:  l.Filename,
		MaxSize:    10,
		MaxBackups: 10,
		//MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (l *Logger) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start).Microseconds()
		l.ReqLog(c, cost)
		c.Next()
	}
}

func (l *Logger) ReqLog(c *gin.Context, cost int64) {
	statusCode := c.Writer.Status()
	clientIp := c.ClientIP()
	userAgent := c.Request.UserAgent()
	dataSize := c.Writer.Size()
	if dataSize < 0 {
		dataSize = 0
	}
	method := c.Request.Method
	path := c.Request.RequestURI

	field := []zap.Field{
		zap.Int("statusCode", statusCode),
		zap.String("ip", clientIp),
		zap.String("userAgent", userAgent),
		zap.Int("dataSize", dataSize),
		zap.String("method", method),
		zap.String("path", path),
		zap.Int64("cost", cost),
	}

	if len(c.Errors) > 0 {
		field = append(field, zap.String("errMsg", c.Errors.ByType(gin.ErrorTypePrivate).String()))
	}
	if statusCode >= 500 {
		l.logger.Error("reqLog", field...)
	} else if statusCode >= 400 {
		l.logger.Warn("reqLog", field...)
	} else {
		l.logger.Info("reqLog", field...)
	}

	l.logger.Info("reqLog", field...)

}