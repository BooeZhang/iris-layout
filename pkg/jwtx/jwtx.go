package jwtx

import (
	"time"

	"github.com/kataras/iris/v12"

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

// VerifyMiddleware token 校验中间件
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
		_ = ctx.StopWithJSON(401, iris.Map{
			"code": 401,
			"msg":  "token 无效",
			"data": nil})
	}
	return verifier.Verify(func() any {
		return new(UserClaims)
	})
}

// GetClaims 获取 token Claims
func GetClaims(ctx iris.Context) *UserClaims {
	claims, ok := jwt.Get(ctx).(*UserClaims)
	if ok {
		return claims
	}
	return &UserClaims{}
}

// GetUserID Get user id
func GetUserID(ctx iris.Context) uint {
	return GetClaims(ctx).UserId
}

// GetUserName 获取用户名
func GetUserName(ctx iris.Context) string {
	return GetClaims(ctx).UserName
}
