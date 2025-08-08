package interfaces

import "blog-api/Domain/models"

type BlogRepository interface {
	CreateBlog(blog models.Blog) (models.Blog, error)
	GetPaginatedBlogs(page, limit int) ([]models.Blog, error)
	GetBlogByID(blogID string) (models.Blog, error)
	UpdateBlog(blog models.Blog) (models.Blog, error)
	DeleteBlog(blogID string) error

	SearchBlogs(query string) ([]models.Blog, error)
	FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]models.Blog, error)

	IncrementViewCount(blogID string) error
	UpdateLikes(blogID string, increment bool) error
	UpdateDislikes(blogID string, increment bool) error
	AddComment(blogID string, comment models.Comment) (models.Comment, error)
	GetComments(blogID string) ([]models.Comment, error)
}
