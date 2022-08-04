package sys_config_dto

import "github.com/unhejing/go-gin-api/model/sys_config_model"

type PageRes struct {
	List       []sys_config_model.SysConfig `json:"list"`
	Pagination struct {
		Total        int `json:"total"`
		CurrentPage  int `json:"current_page"`
		PerPageCount int `json:"per_page_count"`
	} `json:"pagination"`
}
