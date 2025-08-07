package usecases

import (
	"blog-api/Domain/interfaces"
	"blog-api/Domain/models"
)

type BlogUseCase interface {
	CreateBlog(blog models.Blog) (models.Blog, error)
	GetPaginatedBlogs(page, limit int) ([]models.Blog, error)
	GetBlogByID(blogID string) (models.Blog, error)
	UpdateBlog(blog models.Blog) (models.Blog, error)
	DeleteBlog(blogID string) error
	SearchBlogs(query string) ([]models.Blog, error)
	FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]models.Blog, error)

	// popularity tracking methods
	IncrementViewCount(blogID string) error
	UpdateLikes(blogID string, increment bool) error
	UpdateDislikes(blogID string, increment bool) error
	AddComment(blogID string, comment models.Comment) (models.Comment, error)
	GetComments(blogID string) ([]models.Comment, error)
}

type blogUseCase struct {
	blogRepo interfaces.BlogRepository
}

func NewBlogUseCase(blogRepo interfaces.BlogRepository) BlogUseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
	}
}

func (b *blogUseCase) CreateBlog(blog models.Blog) (models.Blog, error) {
	return b.blogRepo.CreateBlog(blog)
}

func (b *blogUseCase) GetPaginatedBlogs(page, limit int) ([]models.Blog, error) {
	return b.blogRepo.GetPaginatedBlogs(page, limit)
}

func (b *blogUseCase) GetBlogByID(blogID string) (models.Blog, error) {
	return b.blogRepo.GetBlogByID(blogID)
}

func (b *blogUseCase) UpdateBlog(blog models.Blog) (models.Blog, error) {
	return b.blogRepo.UpdateBlog(blog)
}

func (b *blogUseCase) DeleteBlog(blogID string) error {
	return b.blogRepo.DeleteBlog(blogID)
}

func (b *blogUseCase) SearchBlogs(query string) ([]models.Blog, error) {
	return b.blogRepo.SearchBlogs(query)
}

func (b *blogUseCase) FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]models.Blog, error) {
	return b.blogRepo.FilterBlogs(tags, dateRange, sortBy)
}

func (b *blogUseCase) IncrementViewCount(blogID string) error {
	return b.blogRepo.IncrementViewCount(blogID)
}

func (b *blogUseCase) UpdateLikes(blogID string, increment bool) error {
	return b.blogRepo.UpdateLikes(blogID, increment)
}

func (b *blogUseCase) UpdateDislikes(blogID string, increment bool) error {
	return b.blogRepo.UpdateDislikes(blogID, increment)
}

func (b *blogUseCase) AddComment(blogID string, comment models.Comment) (models.Comment, error) {
	return b.blogRepo.AddComment(blogID, comment)
}

func (b *blogUseCase) GetComments(blogID string) ([]models.Comment, error) {
	return b.blogRepo.GetComments(blogID)
}
