package response

import (
	"github.com/kataras/iris/v12"

	"irir-layout/config"
	"irir-layout/pkg/erroron"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
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

type validationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

func Error(ctx iris.Context, err error, data any) {
	// var validatorErrs validator.ValidationErrors
	// if errors.As(err, &validatorErrs) {
	// 	errs := make([]validationError, 0, len(validatorErrs))
	// 	for i := range validatorErrs {
	// 		validationErr := validatorErrs[i]
	// 		errs = append(errs, validationError{
	// 			ActualTag: validationErr.ActualTag(),
	// 			Namespace: validationErr.Namespace(),
	// 			Kind:      validationErr.Kind().String(),
	// 			Type:      validationErr.Type().String(),
	// 			Value:     fmt.Sprintf("%v", validationErr.Value()),
	// 			Param:     validationErr.Param(),
	// 		})
	// 	}
	// 	log.Errorf("%+v", errs)
	// 	Ok(c, erroron.ErrParameter, nil)
	// 	return
	// }
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
