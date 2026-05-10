package postgres

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/models"
	"gorm.io/gorm"
)

type pgPost struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"type:varchar(255);not null"`
	Content    string `gorm:"type:text;not null"`
	Author     string `gorm:"type:varchar(100);not null"`
	Extensions string `gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Comments   []pgComment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

type pgComment struct {
	ID               uint    `gorm:"primaryKey"`
	PostID           uint    `gorm:"index;not null"`
	Content          string  `gorm:"type:text;not null"`
	Author           string  `gorm:"type:varchar(100);not null"`
	ParentID         *uint   `gorm:"index"`
	ReplyToAuthor    *string `gorm:"type:varchar(100)"`
	ReplyToCommentID *uint   `gorm:"index"`
	CreatedAt        time.Time
	Children         []pgComment `gorm:"-"`
}

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&pgPost{}, &pgComment{})
}

func (r *Repo) CreatePost(req models.CreatePostReq) (*models.Post, error) {
	extJSON, err := json.Marshal(req.Extensions)
	if err != nil {
		return nil, err
	}
	if req.Extensions == nil {
		extJSON = []byte("{}")
	}

	post := pgPost{
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		Extensions: string(extJSON),
	}
	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return toPost(&post), nil
}

func (r *Repo) GetPost(id string) (*models.Post, error) {
	uid, err := parseUint(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	var post pgPost
	if err := r.db.Preload("Comments").First(&post, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return toPostWithTree(&post), nil
}

func (r *Repo) ListPosts() ([]models.Post, error) {
	var posts []pgPost
	if err := r.db.Preload("Comments").Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	result := make([]models.Post, 0, len(posts))
	for i := range posts {
		result = append(result, *toPostWithTree(&posts[i]))
	}
	return result, nil
}

func (r *Repo) UpdatePost(id string, req models.UpdatePostReq) (*models.Post, error) {
	uid, err := parseUint(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Extensions != nil {
		extJSON, err := json.Marshal(req.Extensions)
		if err != nil {
			return nil, err
		}
		updates["extensions"] = string(extJSON)
	}

	if err := r.db.Model(&pgPost{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetPost(id)
}

func (r *Repo) DeletePost(id string) error {
	uid, err := parseUint(id)
	if err != nil {
		return errors.New("invalid id")
	}
	res := r.db.Delete(&pgPost{}, uid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *Repo) AddComment(postID string, req models.CreateCommentReq) (*models.Comment, error) {
	uid, err := parseUint(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}

	comment := pgComment{
		PostID:        uid,
		Content:       req.Content,
		Author:        req.Author,
		ReplyToAuthor: req.ReplyToAuthor,
	}
	if req.ParentID != nil {
		pid, err := parseUint(*req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent comment id")
		}
		comment.ParentID = &pid
	}
	if req.ReplyToCommentID != nil {
		rid, err := parseUint(*req.ReplyToCommentID)
		if err != nil {
			return nil, errors.New("invalid reply_to_comment_id")
		}
		comment.ReplyToCommentID = &rid
	}

	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	result := toPgComment(&comment)
	result.PostID = postID
	return result, nil
}

func (r *Repo) DeleteComment(commentID string) error {
	uid, err := parseUint(commentID)
	if err != nil {
		return errors.New("invalid comment id")
	}
	res := r.db.Delete(&pgComment{}, uid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func (r *Repo) GetPostWithComments(postID string) (*models.Post, error) {
	return r.GetPost(postID)
}

func parseUint(s string) (uint, error) {
	var v uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("not a number")
		}
		v = v*10 + uint(c-'0')
	}
	return v, nil
}

func toPost(p *pgPost) *models.Post {
	ext := map[string]interface{}{}
	if p.Extensions != "" && p.Extensions != "{}" {
		json.Unmarshal([]byte(p.Extensions), &ext)
	}
	return &models.Post{
		ID:         formatUint(p.ID),
		Title:      p.Title,
		Content:    p.Content,
		Author:     p.Author,
		Extensions: ext,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func formatUint(v uint) string {
	if v == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for v > 0 {
		digits = append(digits, byte('0'+v%10))
		v /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}

func toPgComment(c *pgComment) *models.Comment {
	com := &models.Comment{
		ID:            formatUint(c.ID),
		PostID:        formatUint(c.PostID),
		Content:       c.Content,
		Author:        c.Author,
		ReplyToAuthor: c.ReplyToAuthor,
		CreatedAt:     c.CreatedAt,
	}
	if c.ParentID != nil {
		s := formatUint(*c.ParentID)
		com.ParentID = &s
	}
	if c.ReplyToCommentID != nil {
		s := formatUint(*c.ReplyToCommentID)
		com.ReplyToCommentID = &s
	}
	return com
}

func toPostWithTree(p *pgPost) *models.Post {
	post := toPost(p)
	if len(p.Comments) == 0 {
		post.Comments = []models.Comment{}
		return post
	}

	commentMap := make(map[string]*models.Comment, len(p.Comments))
	for i := range p.Comments {
		c := toPgComment(&p.Comments[i])
		c.Children = []models.Comment{}
		commentMap[c.ID] = c
	}

	var roots []*models.Comment
	for _, c := range commentMap {
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			parent, ok := commentMap[*c.ParentID]
			if ok {
				parent.Children = append(parent.Children, *c)
			} else {
				roots = append(roots, c)
			}
		}
	}

	result := make([]models.Comment, 0, len(roots))
	for _, r := range roots {
		result = append(result, *r)
	}
	post.Comments = result
	return post
}
