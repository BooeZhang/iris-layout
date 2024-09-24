package admin

import (
	"github.com/kataras/iris/v12"

	srvv1 "irir-layout/internal/service/v1/admin"
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
// @Success 200 {object} response.Response{data=schema.LoginRes} "ok"
// @Router /user/login/ [post]
func (cc CommonController) Login(ctx iris.Context) {
	data, err := cc.srv.Login(ctx, "", "")
	if err != nil {
		response.Error(ctx, err, nil)
	}
	response.Ok(ctx, nil, data)
}
