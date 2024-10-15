package router

import (
	"github.com/kataras/iris/v12"

	"irir-layout/pkg/jwtx"
	"irir-layout/pkg/log"

	"irir-layout/internal/controller/v1/admin"
)

type _admin struct{}

var Admin = _admin{}

func (_admin) Load(r *iris.Application) {
	comm := admin.NewCommonController()
	r.Post("/v1/login", comm.Login)
	r.Get("/v1/index", jwtx.VerifyMiddleware(), func(ctx iris.Context) {
		log.Info(ctx, "======")
		_ = ctx.JSON(iris.Map{"status": "ok"})
	})

	r.Use(jwtx.VerifyMiddleware())
}
