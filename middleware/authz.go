package middleware

import (
	"fmt"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/kataras/iris/v12"

	"irir-layout/pkg/jwtx"
	"irir-layout/pkg/log"
)

type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizer(e *casbin.Enforcer) iris.Handler {
	a := &BasicAuthorizer{enforcer: e}
	return func(ctx iris.Context) {
		if !a.CheckPermission(ctx) {
			a.RequirePermission(ctx)
		}
		ctx.Next()
	}
}

func (a *BasicAuthorizer) CheckPermission(ctx iris.Context) bool {
	_ = a.enforcer.LoadPolicy()
	userId := a.GetUserId(ctx)
	method := ctx.Method()
	path := ctx.Path()

	allowed, err := a.enforcer.Enforce(userId, path, method)
	if err != nil {
		ctx.Application().Logger().Error(fmt.Sprintf("校验权限失败: %s", err), log.Fields(ctx))
		return false
	}

	return allowed
}

func (a *BasicAuthorizer) GetUserId(ctx iris.Context) string {
	userId := jwtx.GetUserID(ctx)
	return strconv.FormatInt(int64(userId), 10)
}

func (a *BasicAuthorizer) RequirePermission(ctx iris.Context) {
	ctx.StopWithStatus(iris.StatusForbidden)
}
