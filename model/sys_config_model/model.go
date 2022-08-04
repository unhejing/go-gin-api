package sys_config_model

import "time"

// SysConfig 配置表
//go:generate gormgen -structs SysConfig -input .
type SysConfig struct {
	Id         int32  `json:"id"`          // id
	Note       string `json:"note"`        // 备注
	ParamsName string `json:"params_name"` // 参数名称
	ParamsKey  string `json:"params_key"`  // 参数键名

	ParamsValue string    `json:"params_value"`             // 参数键值
	ChannelTag  string    `json:"channel_tag"`              // 渠道标签
	CreatedTime time.Time `json:"created_time" gorm:"time"` // 创建时间
	UpdatedTime time.Time `json:"updated_time" gorm:"time"` // 更新时间
}
