package usecases

import (
	"blog-api/Domain/interfaces"
	domain "blog-api/Domain/models"
)

// type BlogUseCase interface {
// 	CreateBlog(blog domain.Blog) (domain.Blog, error)
// }

type BlogUseCase interface {
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

type blogUseCase struct {
	blogRepo interfaces.BlogRepository
}

func NewBlogUseCase(blogRepo interfaces.BlogRepository) BlogUseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
	}
}
func (b *blogUseCase) CreateBlog(blog domain.Blog) (domain.Blog, error) {
	return b.blogRepo.CreateBlog(blog)
}

func (b *blogUseCase) GetPaginatedBlogs(page, limit int) ([]domain.Blog, error) {
	return b.blogRepo.GetPaginatedBlogs(page, limit)
}

func (b *blogUseCase) GetBlogByID(blogID string) (domain.Blog, error) {
	return b.blogRepo.GetBlogByID(blogID)
}

func (b *blogUseCase) UpdateBlog(blog domain.Blog) (domain.Blog, error) {
	return b.blogRepo.UpdateBlog((blog))
}

func (b *blogUseCase) DeleteBlog(blogId string) error {
	return b.blogRepo.DeleteBlog(blogId)
}

func (b *blogUseCase) SearchBlogs(query string) ([]domain.Blog, error) {
	return b.blogRepo.SearchBlogs(query)
}

func (b *blogUseCase) FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]domain.Blog, error) {
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

func (b *blogUseCase) AddComment(blogID string, comment domain.Comment) (domain.Comment, error) {
	return b.blogRepo.AddComment(blogID, comment)
}

func (b *blogUseCase) GetComments(blogID string) ([]domain.Comment, error) {
	return b.blogRepo.GetComments(blogID)
}
