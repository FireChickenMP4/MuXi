package models

import "time"

type Post struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Author     string                 `json:"author"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Comments   []Comment              `json:"comments,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type CreatePostReq struct {
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Author     string                 `json:"author"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type UpdatePostReq struct {
	Title      *string                `json:"title,omitempty"`
	Content    *string                `json:"content,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type Comment struct {
	ID               string    `json:"id"`
	PostID           string    `json:"post_id"`
	Content          string    `json:"content"`
	Author           string    `json:"author"`
	ParentID         *string   `json:"parent_id,omitempty"`
	ReplyToAuthor    *string   `json:"reply_to_author,omitempty"`
	ReplyToCommentID *string   `json:"reply_to_comment_id,omitempty"`
	Children         []Comment `json:"children,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateCommentReq struct {
	Content          string  `json:"content"`
	Author           string  `json:"author"`
	ParentID         *string `json:"parent_id,omitempty"`
	ReplyToAuthor    *string `json:"reply_to_author,omitempty"`
	ReplyToCommentID *string `json:"reply_to_comment_id,omitempty"`
}

type Repository interface {
	// Post operations
	CreatePost(req CreatePostReq) (*Post, error)
	GetPost(id string) (*Post, error)
	ListPosts() ([]Post, error)
	UpdatePost(id string, req UpdatePostReq) (*Post, error)
	DeletePost(id string) error

	// Comment operations
	AddComment(postID string, req CreateCommentReq) (*Comment, error)
	DeleteComment(commentID string) error
	GetPostWithComments(postID string) (*Post, error)
}
