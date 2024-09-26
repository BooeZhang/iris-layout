package schema

// LoginReq 登陆请求参数
type LoginReq struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
