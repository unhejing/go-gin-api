package sys_user_dto

type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}
