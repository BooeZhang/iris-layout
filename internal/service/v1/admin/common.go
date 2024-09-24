package admin

import (
	"github.com/kataras/iris/v12"

	v1 "irir-layout/internal/service/v1"
	"irir-layout/pkg/jwtx"
	"irir-layout/pkg/schema"
)

type CommService struct {
	ctx *v1.ServiceContext
}

func NewCommService() *CommService {
	return &CommService{ctx: v1.NewServiceContext()}
}

func (cs *CommService) Login(ctx iris.Context, name, pwd string) (*schema.LoginRes, error) {
	claims := jwtx.UserClaims{
		UserId:   1,
		UserName: "test",
	}

	token, err := jwtx.Signer.Sign(claims)
	if err != nil {
		return nil, err
	}

	return &schema.LoginRes{
		Token: string(token),
	}, nil

}
