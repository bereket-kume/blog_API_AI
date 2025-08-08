// controllers/dto.go
package controllers

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Blog DTOs
type CreateBlogRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
}

type UpdateBlogRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	IsPublished *bool    `json:"is_published"`
}

type BlogResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	AuthorID    string            `json:"author_id"`
	AuthorName  string            `json:"author_name"`
	Tags        []string          `json:"tags"`
	ViewCount   int               `json:"view_count"`
	Likes       int               `json:"likes"`
	Dislikes    int               `json:"dislikes"`
	Comments    []CommentResponse `json:"comments"`
	IsPublished bool              `json:"is_published"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID         string `json:"id"`
	BlogID     string `json:"blog_id"`
	AuthorID   string `json:"author_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PaginatedBlogsResponse struct {
	Blogs      []BlogResponse `json:"blogs"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type SearchBlogsRequest struct {
	Query string `json:"query" binding:"required"`
}

type FilterBlogsRequest struct {
	Tags      []string `json:"tags"`
	DateRange []string `json:"date_range"`
	SortBy    string   `json:"sort_by"`
}
