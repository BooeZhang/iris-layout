package router

import (
	"github.com/kataras/iris/v12"

	"irir-layout/internal/controller/v1/admin"
)

type _admin struct{}

var Admin = _admin{}

func (_admin) Load(r *iris.Application) {
	comm := admin.NewCommonController()
	r.Post("/v1/login", comm.Login)
}
