package dto

type NoticeReq struct {
	ReadType int    `json:"read_type"`
	Title    string `json:"title"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type MarkReadNoticeReq struct {
	Ids []int `json:"ids"`
}
