package interfaces

import (
	domain "blog-api/Domain/models"
)

type BlogRepository interface {
	CreateBlog(blog domain.Blog) (domain.Blog, error)
	GetPaginatedBlogs(page, limit int) ([]domain.Blog, error)
	GetBlogByID(blogID string) (domain.Blog, error)
	UpdateBlog(blog domain.Blog) (domain.Blog, error)
	DeleteBlog(blogID string) error
	SearchBlogs(query string) ([]domain.Blog, error)
	FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]domain.Blog, error)

	// popularity tracking methods
	IncrementViewCount(blogID string) error
	UpdateLikes(blogID string, increment bool) error
	UpdateDislikes(blogID string, increment bool) error
	AddComment(blogID string, comment domain.Comment) (domain.Comment, error)
	GetComments(blogID string) ([]domain.Comment, error)
}
