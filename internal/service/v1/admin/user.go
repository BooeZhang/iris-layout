package admin

import (
	"github.com/kataras/iris/v12"

	"irir-layout/internal/model"
	"irir-layout/internal/repo"
	rp "irir-layout/internal/repo/mysql"
	"irir-layout/pkg/erroron"
	"irir-layout/pkg/jwtx"
	"irir-layout/store/mysql"
)

type UserService struct {
	userRepo repo.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: rp.NewUserRepo(mysql.GetDB()),
	}
}

func (cs UserService) Login(ctx iris.Context, name, pwd string) (model.LoginRes, error) {
	var (
		res model.LoginRes
	)
	user, err := cs.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return res, err
	}

	if user == nil || user.ID == 0 {
		return res, erroron.ErrNotFoundUser
	}

	if user.Compare(pwd) != nil {
		return res, erroron.ErrUserNameOrPwd
	}

	claims := jwtx.UserClaims{
		UserId:   user.ID,
		UserName: user.Account,
	}
	res.AccessToken, err = jwtx.GenAccessToken(claims)
	res.RefreshToken, err = jwtx.GenRefreshToken(claims)
	if err != nil {
		return res, err
	}
	return res, nil

}
