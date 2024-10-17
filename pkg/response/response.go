package response

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"

	"irir-layout/config"
	"irir-layout/pkg/erroron"
	"irir-layout/pkg/log"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type pages struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	List     any   `json:"list"`
}

// Ok 通用响应
func Ok(ctx iris.Context, err error, data any) {
	code, _, msg := erroron.DecodeErr(err)
	r := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	sendErr := ctx.JSON(r, iris.JSON{Indent: "", Secure: true})
	if sendErr != nil {
		ctx.Application().Logger().Warnf("send msg: %s", sendErr)
	}
}

func PageOk(ctx iris.Context, err error, data interface{}, total int64, page, pageSize int) {
	p := pages{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     data,
	}
	Ok(ctx, err, p)
}

type validationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

func Error(ctx iris.Context, err error, data any) {
	var validatorErrs validator.ValidationErrors
	if errors.As(err, &validatorErrs) {
		ctx.Application().Logger().Error(fmt.Sprintf("%s --- %s", ctx.Request().URL, err), log.Fields(ctx))
		validationErrors := wrapValidationErrors(validatorErrs)
		if config.GetConfig().HttpServerConfig.Debug {
			Ok(ctx, erroron.ErrParameter, validationErrors)
			return
		}
		Ok(ctx, erroron.ErrParameter, nil)
		return
	}
	code, httpCode, msg := erroron.DecodeErr(err)
	if !config.GetConfig().HttpServerConfig.Debug && code == 500 {
		msg = "服务器内部错误"
	}

	ctx.StopWithStatus(httpCode)
	sendErr := ctx.StopWithJSON(httpCode, iris.Map{
		"code": code,
		"msg":  msg,
		"data": data})
	if sendErr != nil {
		ctx.Application().Logger().Warnf("send msg: %s", sendErr)
	}
}

func wrapValidationErrors(errs validator.ValidationErrors) []validationError {
	validationErrors := make([]validationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, validationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
		})
	}

	return validationErrors
}
