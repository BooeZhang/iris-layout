package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/requestid"

	"irir-layout/pkg/jwtx"
)

// 获取日志文件名和行号
func getLogSource() (file string, line int) {
	file = "???"
	line = 0
	pc := make([]uintptr, 64)
	n := runtime.Callers(3, pc)
	if n != 0 {
		pc = pc[:n]
		frames := runtime.CallersFrames(pc)

		for {
			frame, more := frames.Next()
			if !strings.Contains(frame.File, "github.com/kataras/golog") {
				file = frame.File
				line = frame.Line
				break
			}
			if !more {
				break
			}
		}
	}

	slices := strings.Split(file, "/")
	file = slices[len(slices)-1]
	return file, line
}

func TextHandler(l *golog.Log) bool {
	file, line := getLogSource()
	l.Message = fmt.Sprintf("[%s:%d] %s", file, line, l.Message)
	return false
}

func JSONHandler(l *golog.Log) bool {
	file, line := getLogSource()
	l.Fields = golog.Fields{"file": file, "line": line}
	return false
}

func FuncCtx(ctx iris.Context, latency time.Duration) {
	var status, ip, method, path, _requestID string
	ip = ctx.RemoteAddr()
	method = ctx.Method()
	path = ctx.Request().URL.Path
	status = strconv.Itoa(ctx.GetStatusCode())
	_requestID = requestid.Get(ctx)

	line := fmt.Sprintf("%s %4v %s %v %s", status, latency, ip, method, path)

	fields := golog.Fields{"requestID": _requestID}
	userName := jwtx.GetUserName(ctx)
	if len(userName) > 0 {
		fields["userName"] = userName
	}
	if context.StatusCodeNotSuccessful(ctx.GetStatusCode()) {
		ctx.Application().Logger().Warn(line, fields)
	} else {
		ctx.Application().Logger().Info(line, fields)
	}
}

func Fields(ctx iris.Context) golog.Fields {
	return golog.Fields{"requestID": requestid.Get(ctx)}
}

func Debug(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.DebugLevel, msg, Fields(ctx))
}

func Debugf(ctx iris.Context, format string, args ...interface{}) {
	args = append(args, Fields(ctx))
	ctx.Application().Logger().Logf(golog.DebugLevel, format, args)
}

func Info(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.InfoLevel, msg, Fields(ctx))
}

func Infof(ctx iris.Context, format string, args ...any) {
	args = append(args, Fields(ctx))
	ctx.Application().Logger().Logf(golog.InfoLevel, format, args)
}

func Warn(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.WarnLevel, msg, Fields(ctx))
}
func Warnf(ctx iris.Context, format string, args ...any) {
	args = append(args, Fields(ctx))
	ctx.Application().Logger().Logf(golog.WarnLevel, format, args)
}

func Error(ctx iris.Context, msg any) {
	ctx.Application().Logger().Log(golog.ErrorLevel, msg, Fields(ctx))
}

func Errorf(ctx iris.Context, format string, args ...any) {
	args = append(args, Fields(ctx))
	ctx.Application().Logger().Logf(golog.ErrorLevel, format, args...)
}
