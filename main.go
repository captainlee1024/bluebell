// Package main provides ...
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/captainlee1024/bluebell/dao/mysql"
	"github.com/captainlee1024/bluebell/dao/redis"
	"github.com/captainlee1024/bluebell/logger"
	"github.com/captainlee1024/bluebell/pkg/snowflake"
	"github.com/captainlee1024/bluebell/router"
	"github.com/captainlee1024/bluebell/settings"
	"go.uber.org/zap"
)

/* swagger main 函数注释格式（写项目相关介绍信息）
// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
*/

// @title bluebell
// @version 1.0
// @description
// @tremsOfService http://swagger.io/terms/

// @contact.name CaptainLee1024
// @contact.url http://blog.leecoding.club
// @contact.email 644052732@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/license/LICENSE-2.0html

// @host localhost:8081
// @BasePath /api/v1
func main() {
	// 1. 初始化配
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2. 初始化 logger
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 3. 初始化 MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Error("init mysql failed", zap.Error(err))
		return
	}
	defer mysql.Close()

	// 4. 初始化 Reids
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Error("init redis failed", zap.Error(err))
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		zap.L().Error("init snowflake failed", zap.Error(err))
		return
	}
	// 5. 注册路由
	r := router.Setup(settings.Conf.Mode)

	// 6. 启动程序（优雅关机）
	srv := &http.Server{
		//Addr: fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅关闭服务器，为关闭服务器操作设置一个5秒的延时
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在这里，当收到上述两种信号时才往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
