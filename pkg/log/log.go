package log

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/requestid"
)

func FuncCtx(ctx iris.Context, latency time.Duration) {
	var status, ip, method, path, _requestid string
	ip = ctx.RemoteAddr()
	method = ctx.Method()
	path = ctx.Request().URL.Path
	status = strconv.Itoa(ctx.GetStatusCode())
	_requestid = requestid.Get(ctx)

	line := fmt.Sprintf("%s %v %4v %s %s %s", _requestid, status, latency, ip, method, path)

	if context.StatusCodeNotSuccessful(ctx.GetStatusCode()) {
		ctx.Application().Logger().Warn(line)
	} else {
		ctx.Application().Logger().Info(line)
	}
}

func fields(ctx iris.Context) golog.Fields {
	return golog.Fields{"requestId": requestid.Get(ctx)}
}

func Debug(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.DebugLevel, msg, fields(ctx))
}

func Debugf(ctx iris.Context, format string, args ...interface{}) {
	ctx.Application().Logger().Logf(golog.DebugLevel, format, args, fields(ctx))
}

func Info(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.InfoLevel, msg, fields(ctx))
}

func Infof(ctx iris.Context, format string, args ...any) {
	ctx.Application().Logger().Logf(golog.InfoLevel, format, args, fields(ctx))
}

func Warn(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.WarnLevel, msg, fields(ctx))
}
func Warnf(ctx iris.Context, format string, args ...any) {
	ctx.Application().Logger().Logf(golog.WarnLevel, format, args, fields(ctx))
}

func Error(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.ErrorLevel, msg, fields(ctx))
}

func Errorf(ctx iris.Context, format string, args ...any) {
	ctx.Application().Logger().Logf(golog.ErrorLevel, format, args, fields(ctx))
}
