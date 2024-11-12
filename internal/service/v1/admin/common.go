package admin

import (
	"github.com/kataras/iris/v12"

	"irir-layout/internal/model"
)

type CommService struct {
	urc *UserService
}

func NewCommService() *CommService {
	return &CommService{
		urc: NewUserService(),
	}
}

func (cs CommService) Login(ctx iris.Context, name, pwd string) (model.LoginRes, error) {
	return cs.urc.Login(ctx, name, pwd)
}
