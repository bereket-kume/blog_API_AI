package models

import (
	"time"
)

// Blog represents a blog post in the system
// Following clean architecture: domain layer should be independent of infrastructure
type Blog struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title" bson:"title"`
	Content     string    `json:"content" bson:"content"`
	AuthorID    string    `json:"author_id" bson:"author_id"`
	AuthorName  string    `json:"author_name" bson:"author_name"`
	Tags        []string  `json:"tags" bson:"tags"`
	ViewCount   int       `json:"view_count" bson:"view_count"`
	Likes       int       `json:"likes" bson:"likes"`
	Dislikes    int       `json:"dislikes" bson:"dislikes"`
	Comments    []Comment `json:"comments" bson:"comments"`
	IsPublished bool      `json:"is_published" bson:"is_published"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Comment represents a comment on a blog post
type Comment struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	BlogID     string    `json:"blog_id" bson:"blog_id"`
	AuthorID   string    `json:"author_id" bson:"author_id"`
	AuthorName string    `json:"author_name" bson:"author_name"`
	Content    string    `json:"content" bson:"content"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}
