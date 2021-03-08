package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/towithyou/web_app/dao/mysql"
	"github.com/towithyou/web_app/dao/redis"
	"github.com/towithyou/web_app/logger"
	"github.com/towithyou/web_app/routes"
	"github.com/towithyou/web_app/settings"
	"go.uber.org/zap"
)

// GO web 脚手架

func main() {
	var file string
	flag.StringVar(&file, "config", "config.yaml", "配置文件路径")
	flag.Parse()
	// 加载配置文件
	if err := settings.Init(file); err != nil {
		fmt.Printf("init settings error %v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger error %v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	// 初始化mysql redis
	if err := mysql.Init(settings.Conf.MySqlConfig); err != nil {
		fmt.Printf("init mysql error %v\n", err)
		return
	}
	defer mysql.Close()

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis error %v\n", err)
		return
	}
	defer redis.Close()

	// 注册路由
	r := routes.Setup()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	// 启动服务
	go func() {
		// 开启一个goroutine启动服务
		fmt.Printf("listen port at: %d\n", settings.Conf.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			zap.L().Error("listen error ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")

	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Shutdown Server ...", zap.Error(err))
	}

	zap.L().Info("Server exiting ...")
}
