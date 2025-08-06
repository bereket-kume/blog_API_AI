package repositories

import (
	domain "blog-api/Domain/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	BlogCollection *mongo.Collection
}

func NewBlogRepository(db *mongo.Collection) *BlogRepository {
	return &BlogRepository{
		BlogCollection: db,
	}
}

func (r *BlogRepository) CreateBlog(blog domain.Blog) (domain.Blog, error) {
	result, err := r.BlogCollection.InsertOne(context.Background(), blog)
	if err != nil {
		log.Printf("Error inserting blog into MongoDB: %v", err)
		return domain.Blog{}, err
	}
	blog.ID = result.InsertedID.(primitive.ObjectID)
	return blog, nil
}

func (r *BlogRepository) GetPaginatedBlogs(page, limit int) ([]domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	skip := (page - 1) * limit
	cursor, err := r.BlogCollection.Find(ctx, nil, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		log.Printf("Error fetching paginated blogs: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var blogs []domain.Blog
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			log.Printf("Error decoding blog: %v", err)
			continue
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *BlogRepository) GetBlogByID(blogID string) (domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		return domain.Blog{}, err
	}

	filter := bson.M{"_id": id}
	var blog domain.Blog
	err = r.BlogCollection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Blog not found with ID: %s", blogID)
			return domain.Blog{}, err
		}
		log.Printf("Error fetching blog by ID: %v", err)
		return domain.Blog{}, err
	}

	return blog, nil
}

func (r *BlogRepository) UpdateBlog(blog domain.Blog) (domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": blog.ID}
	update := bson.M{
		"$set": bson.M{
			"title":      blog.Title,
			"content":    blog.Content,
			"tags":       blog.Tags,
			"updated_at": time.Now(),
		},
	}
	_, err := r.BlogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating blog: %v", err)
	}

	return domain.Blog{}, nil
}

func (r *BlogRepository) DeleteBlog(blogID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		return err
	}
	filter := bson.M{"_id": id}
	result, err := r.BlogCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting blog: %v", err)
		return err
	}
	if result.DeletedCount == 0 {
		log.Printf("Blog not found with ID: %s", blogID)
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *BlogRepository) SearchBlogs(query string) ([]domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
			{"tags": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	cursor, err := r.BlogCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error searching blogs: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var blogs []domain.Blog
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			log.Printf("Error decoding blog: %v", err)
			continue
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *BlogRepository) FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]domain.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}
	if dateRange[0] != "" && dateRange[1] != "" {
		startDate, err := time.Parse(time.RFC3339, dateRange[0])
		if err != nil {
			log.Printf("Error parsing start date: %v", err)
			return nil, err
		}
		endDate, err := time.Parse(time.RFC3339, dateRange[1])
		if err != nil {
			log.Printf("Error parsing end date: %v", err)
			return nil, err
		}
		filter["created_at"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}
	options := options.Find()
	if sortBy != "" {
		options.SetSort(bson.M{sortBy: -1})
	}

	cursor, err := r.BlogCollection.Find(ctx, filter, options)
	if err != nil {
		log.Printf("Error filtering blogs: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var blogs []domain.Blog
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			log.Printf("Error decoding blog: %v", err)
			continue
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (r *BlogRepository) IncrementViewCount(blogID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"popularity.views": 1,
		},
	}
	_, err = r.BlogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error incrementing view count: %v", err)
	}

	return nil
}

func (r *BlogRepository) UpdateLikes(blogID string, increment bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"popularity.likes": 1,
		},
	}
	if !increment {
		update = bson.M{
			"$inc": bson.M{
				"popularity.likes": -1,
			},
		}
	}
	_, err = r.BlogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating likes: %v", err)
	}

	return nil
}

func (r *BlogRepository) UpdateDislikes(blogID string, increment bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"popularity.dislikes": 1,
		},
	}
	if !increment {
		update = bson.M{
			"$inc": bson.M{
				"popularity.dislikes": -1,
			},
		}
	}
	_, err = r.BlogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating dislikes: %v", err)
		return err
	}

	return nil
}

func (r *BlogRepository) AddComment(blogID string, comment domain.Comment) (domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	comment.ID = primitive.NewObjectID()
	comment.BlogID, _ = primitive.ObjectIDFromHex(blogID)
	comment.CreatedAt = time.Now()
	_, err := r.BlogCollection.InsertOne(ctx, comment)
	if err != nil {
		log.Printf("Error adding comment: %v", err)
		return domain.Comment{}, err
	}
	// increment the comment count in the blog's popularity
	filter := bson.M{"_id": comment.BlogID}
	update := bson.M{
		"$inc": bson.M{
			"popularity.comments": 1,
		},
	}
	_, err = r.BlogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error incrementing comment count: %v", err)
		return domain.Comment{}, err
	}

	return domain.Comment{}, nil
}

func (r *BlogRepository) GetComments(blogID string) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		return nil, err
	}
	filter := bson.M{"blog_id": id}
	cursor, err := r.BlogCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var comments []domain.Comment
	for cursor.Next(ctx) {
		var comment domain.Comment
		if err := cursor.Decode(&comment); err != nil {
			log.Printf("Error decoding comment: %v", err)
			continue
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
