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

type Join2Req struct {
	Token   string `json:"token"`
	UnionID string `json:"union_id"`
}

type LeaveReq struct {
	Token   string `json:"token"`
	UnionID string `json:"union_id"`
}
