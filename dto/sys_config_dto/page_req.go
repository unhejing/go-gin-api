package sys_config_dto

type PageReq struct {
	Page int `json:"page"` // 第几页
	Size int `json:"size"` // 每页显示条数
}
