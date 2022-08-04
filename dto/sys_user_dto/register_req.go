package sys_user_dto

type RegisterReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Nickname string `json:"nickname" binding:"required"` // 昵称
	Mobile   string `json:"mobile" binding:"required"`   // 手机号
	Password string `json:"password" binding:"required"` // MD5后的密码
}
