package jwtx

import (
	"github.com/kataras/iris/v12"
	"irir-layout/pkg/erroron"
	"irir-layout/pkg/response"
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"

	"irir-layout/config"
)

type UserClaims struct {
	jwt.Claims
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
}

// GenAccessToken 生成访问 token
func GenAccessToken(claims UserClaims) (string, error) {
	cf := config.GetConfig()
	sigKey := []byte(cf.JwtConfig.Key)
	signer := jwt.NewSigner(jwt.HS256, sigKey, cf.JwtConfig.AccessExpired*time.Minute).WithEncryption([]byte(cf.JwtConfig.Salt), nil)
	token, err := signer.Sign(claims)
	return string(token), err
}

// GenRefreshToken 生成刷新 token
func GenRefreshToken(claims UserClaims) (string, error) {
	cf := config.GetConfig()
	sigKey := []byte(cf.JwtConfig.Key)
	signer := jwt.NewSigner(jwt.HS256, sigKey, cf.JwtConfig.RefreshExpired*time.Hour).WithEncryption([]byte(cf.JwtConfig.Salt), nil)
	token, err := signer.Sign(claims)
	return string(token), err
}

func VerifyMiddleware() iris.Handler {
	cf := config.GetConfig()
	sigKey := []byte(cf.JwtConfig.Key)
	verifier := jwt.NewVerifier(jwt.HS256, sigKey).WithDecryption([]byte(cf.JwtConfig.Salt), nil)
	verifier.ErrorHandler = func(ctx iris.Context, err error) {
		if err == nil {
			return
		}
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		response.Error(ctx, erroron.ErrTokenInvalid, nil)
	}
	return verifier.Verify(func() any {
		return new(UserClaims)
	})
}
