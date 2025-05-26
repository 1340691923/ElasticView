package dto

type AddCommentRequest struct {
	PluginID int    `json:"plugin_id"`
	Content  string `json:"content"`
	ParentID int    `json:"parent_id"`
}

type DeleteCommentRequest struct {
	CommentID int `json:"comment_id"`
}

type BlockUserRequest struct {
	UserID int    `json:"user_id"`
	Reason string `json:"reason"`
}

type UnblockUserRequest struct {
	UserID int `json:"user_id"`
}

type ListCommentsRequest struct {
	PluginID int `json:"plugin_id"`
}

type LikeCommentRequest struct {
	CommentID int `json:"comment_id"`
	State     int `json:"state"`
}
