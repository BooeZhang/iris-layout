package authz

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/kataras/golog"

	"irir-layout/config"
	"irir-layout/store/mysql"
)

var (
	e *casbin.Enforcer
)

func InitAuth() {
	cf := config.GetConfig()
	a, err := gormadapter.NewAdapterByDB(mysql.GetDB())
	if err != nil {
		golog.Fatal("初始化权限访问数据库失败")
	}

	e, err = casbin.NewEnforcer(cf.HttpServerConfig.CasbinModelPath, a)
	if err != nil {
		golog.Fatal("初始化权限访问系统执行器失败")
	}
}

func GetEnforcer() *casbin.Enforcer {
	return e
}
