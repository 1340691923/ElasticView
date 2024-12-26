package dto

import "github.com/1340691923/ElasticView/pkg/services/gm_user"

type GetMetaByPathReq struct {
	FullPath string           `json:"full_path"`
	Routers  []*gm_user.Route `json:"routers"`
}
