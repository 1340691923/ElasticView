package vo

type Comment struct {
	ID        int64      `db:"id" json:"id"`
	PluginID  int64      `db:"plugin_id" json:"plugin_id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	Content   string     `db:"content" json:"content"`
	ParentID  int64      `db:"parent_id" json:"parent_id"`
	LikeCount int        `db:"like_count" json:"like_count"`
	CreatedAt string     `db:"created_at" json:"created_at"`
	Nickname  string     `db:"realname" json:"realname"`
	HasLike   bool       `json:"has_like" db:"-"`
	Children  []*Comment `json:"children,omitempty"`
}
