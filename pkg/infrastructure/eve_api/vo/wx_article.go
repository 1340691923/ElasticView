package vo

type WxArticleModel struct {
	ID         int    `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Link       string `json:"link" db:"link"`
	CreateTime string `json:"create_time" db:"create_time"`
	Typ        string `json:"typ" db:"typ"` // info/warning/success/danger/primary
	TagName    string `json:"tag_name" db:"tag_name"`
}

type WxArticleList struct {
	Data []WxArticleModel
}
