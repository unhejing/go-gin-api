package sys_user_dto

type SearchOneData struct {
	Id       int32  `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Mobile   string `json:"mobile"`   // 手机号
	Password string `json:"password"` // 密码
	IsUsed   int32  `json:"is_used"`  // 是否启用 1:是  -1:否
}
