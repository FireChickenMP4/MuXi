package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoComment struct {
	ID               primitive.ObjectID  `bson:"_id"`
	Content          string              `bson:"content"`
	Author           string              `bson:"author"`
	ParentID         *primitive.ObjectID `bson:"parent_id,omitempty"`
	ReplyToAuthor    *string             `bson:"reply_to_author,omitempty"`
	ReplyToCommentID *primitive.ObjectID `bson:"reply_to_comment_id,omitempty"`
	CreatedAt        time.Time           `bson:"created_at"`
}

type mongoPost struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty"`
	Title      string                 `bson:"title"`
	Content    string                 `bson:"content"`
	Author     string                 `bson:"author"`
	Extensions map[string]interface{} `bson:"extensions,omitempty"`
	Comments   []mongoComment         `bson:"comments,omitempty"`
	CreatedAt  time.Time              `bson:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at"`
}

type Repo struct {
	col *mongo.Collection
}

func NewRepo(client *mongo.Client, dbName string) *Repo {
	return &Repo{col: client.Database(dbName).Collection("posts")}
}

func (r *Repo) CreatePost(req models.CreatePostReq) (*models.Post, error) {
	now := time.Now()
	doc := mongoPost{
		ID:         primitive.NewObjectID(),
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		Extensions: req.Extensions,
		Comments:   []mongoComment{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if doc.Extensions == nil {
		doc.Extensions = map[string]interface{}{}
	}
	_, err := r.col.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, err
	}
	return toPost(&doc), nil
}

func (r *Repo) GetPost(id string) (*models.Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	var doc mongoPost
	err = r.col.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return toPostWithTree(&doc), nil
}

func (r *Repo) ListPosts() ([]models.Post, error) {
	cursor, err := r.col.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var docs []mongoPost
	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0, len(docs))
	for i := range docs {
		posts = append(posts, *toPostWithTree(&docs[i]))
	}
	return posts, nil
}

func (r *Repo) UpdatePost(id string, req models.UpdatePostReq) (*models.Post, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	update := bson.M{"updated_at": time.Now()}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Content != nil {
		update["content"] = *req.Content
	}
	if req.Extensions != nil {
		update["extensions"] = req.Extensions
	}

	_, err = r.col.UpdateOne(context.Background(), bson.M{"_id": oid}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return r.GetPost(id)
}

func (r *Repo) DeletePost(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	res, err := r.col.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *Repo) AddComment(postID string, req models.CreateCommentReq) (*models.Comment, error) {
	oid, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.New("invalid post id")
	}

	comment := mongoComment{
		ID:        primitive.NewObjectID(),
		Content:   req.Content,
		Author:    req.Author,
		CreatedAt: time.Now(),
	}
	if req.ParentID != nil {
		pid, err := primitive.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent comment id")
		}
		comment.ParentID = &pid
	}
	comment.ReplyToAuthor = req.ReplyToAuthor
	if req.ReplyToCommentID != nil {
		rid, err := primitive.ObjectIDFromHex(*req.ReplyToCommentID)
		if err != nil {
			return nil, errors.New("invalid reply_to_comment_id")
		}
		comment.ReplyToCommentID = &rid
	}

	_, err = r.col.UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$push": bson.M{"comments": comment}},
	)
	if err != nil {
		return nil, err
	}

	result := toComment(&comment)
	result.PostID = postID
	return result, nil
}

func (r *Repo) DeleteComment(commentID string) error {
	oid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New("invalid comment id")
	}

	res, err := r.col.UpdateOne(
		context.Background(),
		bson.M{"comments._id": oid},
		bson.M{"$pull": bson.M{"comments": bson.M{"_id": oid}}},
	)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func (r *Repo) GetPostWithComments(postID string) (*models.Post, error) {
	return r.GetPost(postID)
}

func toPost(doc *mongoPost) *models.Post {
	p := &models.Post{
		ID:         doc.ID.Hex(),
		Title:      doc.Title,
		Content:    doc.Content,
		Author:     doc.Author,
		Extensions: doc.Extensions,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}
	if p.Extensions == nil {
		p.Extensions = map[string]interface{}{}
	}
	return p
}

func toComment(c *mongoComment) *models.Comment {
	com := &models.Comment{
		ID:            c.ID.Hex(),
		Content:       c.Content,
		Author:        c.Author,
		ReplyToAuthor: c.ReplyToAuthor,
		CreatedAt:     c.CreatedAt,
	}
	if c.ParentID != nil {
		s := c.ParentID.Hex()
		com.ParentID = &s
	}
	if c.ReplyToCommentID != nil {
		s := c.ReplyToCommentID.Hex()
		com.ReplyToCommentID = &s
	}
	return com
}

func toPostWithTree(doc *mongoPost) *models.Post {
	p := toPost(doc)
	if len(doc.Comments) == 0 {
		p.Comments = []models.Comment{}
		return p
	}

	commentMap := make(map[string]*models.Comment, len(doc.Comments))
	for i := range doc.Comments {
		c := toComment(&doc.Comments[i])
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
	p.Comments = result
	return p
}
