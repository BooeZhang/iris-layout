package jwtx

import (
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"

	"irir-layout/config"
)

type UserClaims struct {
	UserId   uint
	UserName string
	Expire   time.Time
}

var cf = config.GetConfig()
var sigKey = []byte(cf.JwtConfig.Key)
var Signer = jwt.NewSigner(jwt.HS256, sigKey, cf.JwtConfig.Expired*time.Minute)
var verifier = jwt.NewVerifier(jwt.HS256, sigKey)
var verifyMiddleware = verifier.Verify(func() any { return new(UserClaims) })
