package service

type CreateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"passwrod"`
}

type SignInReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type SignInResp struct {
	Token     string `json:"tocken"`
	Refresh   string `json:"refresh"`
	TokenType string `json:"type"`
	ExpiredIn int64  `json:"exp"`
}
