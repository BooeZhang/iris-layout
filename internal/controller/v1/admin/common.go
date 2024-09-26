package admin

import (
	"github.com/kataras/iris/v12"

	srvv1 "irir-layout/internal/service/v1/admin"
	"irir-layout/pkg/erroron"
	"irir-layout/pkg/response"
	"irir-layout/pkg/schema"
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
// @Success 200 {object} response.Response{data=schema.LoginRes} "ok"
// @Router /login/ [post]
func (cc CommonController) Login(ctx iris.Context) {
	var param schema.LoginReq

	err := ctx.ReadJSON(&param)
	if err != nil {
		ctx.Application().Logger().Error(err.Error())
		response.Error(ctx, erroron.ErrParameter, nil)
		return
	}
	data, err := cc.srv.Login(ctx, param.UserName, param.Password)
	if err != nil {
		response.Error(ctx, err, nil)
		return
	}
	response.Ok(ctx, nil, data)
}
