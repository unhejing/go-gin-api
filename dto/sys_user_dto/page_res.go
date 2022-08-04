package sys_user_dto

import (
	"github.com/unhejing/go-gin-api/model/sys_user_model"
)

type PageRes struct {
	List       []sys_user_model.SysUser `json:"list"`
	Pagination struct {
		Total        int `json:"total"`
		CurrentPage  int `json:"current_page"`
		PerPageCount int `json:"per_page_count"`
	} `json:"pagination"`
}
