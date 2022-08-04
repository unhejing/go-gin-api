package sys_user_model

import "time"

// SysUser 用户表
//go:generate gormgen -structs SysUser -input .
type SysUser struct {
	Id          int32     `json:"id"`                       // 主键
	Username    string    `json:"username"`                 // 用户名
	Password    string    `json:"password"`                 // 密码
	Nickname    string    `json:"nickname"`                 // 昵称
	Mobile      string    `json:"mobile"`                   // 手机号
	IsUsed      int32     `json:"is_used"`                  // 是否启用 1:是  -1:否
	IsDeleted   int32     `json:"is_deleted"`               // 是否删除 1:是  -1:否
	CreatedTime time.Time `json:"created_time" gorm:"time"` // 创建时间
	UpdatedTime time.Time `json:"updated_time" gorm:"time"` // 更新时间
}
