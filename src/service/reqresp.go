package service

type CreateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type SignInResp struct {
	Token     string `json:"token"`
	Refresh   string `json:"refresh"`
	TokenType string `json:"type"`
}
