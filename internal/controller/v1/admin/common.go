package admin

import (
	"fmt"

	"github.com/kataras/iris/v12"

	"irir-layout/internal/model"
	srvv1 "irir-layout/internal/service/v1/admin"
	"irir-layout/pkg/log"
	"irir-layout/pkg/response"
)

type CommonController struct {
	srv *srvv1.CommService
}

func NewCommonController() *CommonController {
	return &CommonController{srv: srvv1.NewCommService()}
}

// Login
// @Summary 登录
// @Schemes
// @Description 登录
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param   data body schema.LoginReq true "."
// @Success 200 {object} response.Response{data=model.LoginRes} "ok"
// @Router /login/ [post]
func (cc CommonController) Login(ctx iris.Context) {
	var param model.LoginReq
	ctx.Application().Logger().Info(fmt.Sprintf("param: %v", param), log.Fields(ctx))
	err := ctx.ReadJSON(&param)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	data, err := cc.srv.Login(ctx, param.UserName, param.Password)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	response.Ok(ctx, nil, data)
}
