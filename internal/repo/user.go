package repo

import (
	"github.com/kataras/iris/v12"

	"irir-layout/internal/model"
)

// UserRepo 用户repo接口
type UserRepo interface {
	GetUserByName(ctx iris.Context, account string) (*model.User, error)
	GetUserById(ctx iris.Context, uid int64) (*model.User, error)
	GetUserByMobile(ctx iris.Context, mobile string) (*model.User, error)
}
