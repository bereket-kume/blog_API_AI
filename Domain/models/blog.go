package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	AuthorID   primitive.ObjectID `bson:"author_id"`
	Tags       []string           `bson:"tags"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at,omitempty"`
	Popularity Popularity         `bson:"popularity"`
}

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at,omitempty"`
}

type Popularity struct {
	Views    int `bson:"views"`
	Likes    int `bson:"likes"`
	Dislikes int `bson:"dislikes"`
	Comments int `bson:"comments"`
}
