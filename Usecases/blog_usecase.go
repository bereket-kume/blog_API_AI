package usecases

import (
	"blog-api/Domain/interfaces"
	domain "blog-api/Domain/models"
)

type BlogUseCase interface {
	CreateBlog(blog domain.Blog) (domain.Blog, error)
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
