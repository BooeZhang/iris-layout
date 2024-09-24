package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"irir-layout/config"
	"irir-layout/core"
	"irir-layout/internal/router"
)

// printWorkingDir 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	golog.Infof("==> 工作目录: %s", wd)
	golog.Infof("==> 使用的配置文件为: `%s`", viper.ConfigFileUsed())
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v2
func main() {
	configFile := flag.String("c", "etc/config.toml", "-c 选项用于指定要使用的配置文件")
	flag.Parse()
	config.InitConfig(*configFile)

	printWorkingDir()
	cf := config.GetConfig()
	app := core.NewHttpServer(cf)
	app.LoadRouter(router.Admin)
	go app.Run()
	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		golog.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if err := app.Application.Shutdown(context.Background()); err != nil {
				golog.Errorf("Server forced to shutdown: %s", err.Error())
			}

			// rpcSrv.GracefulStop()
			// st.Close()
			golog.Info("Server exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
