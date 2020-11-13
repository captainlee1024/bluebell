// Package logger provides ...
package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/captainlee1024/bluebell/settings"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 初始化 Zap logger
// 配置自己的 logger　并替换 zap 中定义的全局变量 logger
func Init(cfg *settings.LogConfig, mode string) (err error) {
	writerSyncer := getLogWriter(
		//viper.GetString("log.filename"),
		//viper.GetInt("log.max_size"),
		//viper.GetInt("log.max_backup"),
		//viper.GetInt("log.max_age"),
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)

	encoder := getEncoder()

	// 把 yaml 配置文件中的 string 类型的 level 配置，解析成 zap 中的 level 类型
	var l = new(zapcore.Level)
	//err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	err = l.UnmarshalText([]byte(viper.GetString(cfg.Level)))
	if err != nil {
		return
	}

	var core zapcore.Core
	if mode == "dev" {
		// 开发模式，输出日志到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee( // 指定两个输出位置
			// 第一个输出　和下面的配置一样，以 JSON 方式写入到日志文件里面
			zapcore.NewCore(encoder, writerSyncer, l),
			// 第二个输出
			// consoleEncoder 指定console编码器
			// zapcore.Lock(os.Stdout) 指定输出位置是标准输出，给它转换成符合条件的 WriteSyncer
			// zapcore.DebugLevel 指定输出日志级别
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)

	} else {
		//
		core = zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	}

	// 生成配置的 logger
	lg := zap.New(core, zap.AddCaller())

	// 替换 zap 中的全局变量
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger", // 名字是什么
		CallerKey:     "caller", // 调用者的名字
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		//EncodeTime: zapcore.EpochTimeEncoder, // 默认的时间编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder, // 修改之后的时间编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			//logger
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
