package main

import (
	"context"
	"fmt"
	cmd "github.com/haleyrom/wallet/cmd/core"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/pkg/middleware"
	"github.com/haleyrom/wallet/pkg/version"
	"github.com/haleyrom/wallet/router"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	// srv srv
	srv *http.Server
)

// init init
func init() {
	cmd.Init()
}

// main main
// @title 测试
// @version 0.0.1
// @description  测试
// @BasePath /v1
func main() {
	// 注册路由
	r := router.InitRouter()
	// 日志中间件
	r.Use(middleware.LoggerToFile())
	_, out := middleware.OpenLoggerFile()
	// 日志文件落地
	logrus.SetOutput(out)
	//设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 版本信息
	version.LogAppInfo()
	srv = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if len(viper.GetString("httpport")) > core.DefaultNilNum {
		srv.Addr = viper.GetString("httpport")
	}
	defer clone()

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// 服务注册
	//middleware.ConsulRegister()
	fmt.Printf("Listening and serving HTTP on %s\n", srv.Addr)

}

// clone 退出
func clone() {
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		_ = core.Orm.Close()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Println("Server exiting")
}
