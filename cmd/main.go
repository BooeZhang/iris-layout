package main

import (
	"flag"
	"os"

	"github.com/kataras/golog"
	"github.com/spf13/viper"

	"irir-layout/config"
	"irir-layout/core"
	"irir-layout/internal/model"
	"irir-layout/internal/router"
	"irir-layout/pkg/authz"
	"irir-layout/pkg/log"
	"irir-layout/store/mysql"
	"irir-layout/store/redis"
)

// printWorkingDir 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	golog.Infof("==> 工作目录: %s", wd)
	golog.Infof("==> 使用的配置文件为: `%s`", viper.ConfigFileUsed())
}

func main() {
	configFile := flag.String("c", "etc/config.toml", "-c 选项用于指定要使用的配置文件")
	flag.Parse()
	config.InitConfig(*configFile)

	printWorkingDir()
	cf := config.GetConfig()

	log.Init(cf)

	mysql.DialToMysql(cf.MysqlConfig)
	defer mysql.Close()
	if cf.HttpServerConfig.Debug {
		migrateDB()
		mysql.CreateSuperUser(mysql.GetDB(), cf.MysqlConfig)
	}

	authz.InitAuth()

	redis.DialToRedis(cf.RedisConfig)
	defer redis.Close()

	http := core.NewHttpServer(cf)
	http.LoadRouter(router.Admin)
	http.Run()
}

func migrateDB() {
	if err := mysql.GetDB().AutoMigrate(
		new(model.User),
	); err != nil {
		golog.Errorf("migrate db failed: %s", err)
		os.Exit(1)
	}
	golog.Info("migrate db completed...")
}
