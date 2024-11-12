package repo

import (
	"github.com/kataras/iris/v12"

	"irir-layout/internal/model"
)

/*
repo 接口相关方法在只在自己的 service 中调用，尽量不要跨 service 使用，
其他 service 如果需要 repo 数据可通过调用相关 service 方法不直接调用 repo.
*/

// UserRepo 用户repo接口
type UserRepo interface {
	GetUserByName(ctx iris.Context, account string) (*model.User, error)
	GetUserById(ctx iris.Context, uid int64) (*model.User, error)
	GetUserByMobile(ctx iris.Context, mobile string) (*model.User, error)
}
