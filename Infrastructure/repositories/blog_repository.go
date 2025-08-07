package repositories

import (
	"blog-api/Domain/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogMongoRepo struct {
	collection *mongo.Collection
}

func NewBlogMongoRepo(col *mongo.Collection) *blogMongoRepo {
	return &blogMongoRepo{collection: col}
}

// CreateBlog creates a new blog post
func (br *blogMongoRepo) CreateBlog(blog models.Blog) (models.Blog, error) {
	blog.ID = primitive.NewObjectID()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	_, err := br.collection.InsertOne(context.TODO(), blog)
	if err != nil {
		return models.Blog{}, err
	}

	return blog, nil
}

// GetPaginatedBlogs retrieves blogs with pagination
func (br *blogMongoRepo) GetPaginatedBlogs(page, limit int) ([]models.Blog, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(skip)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := br.collection.Find(context.TODO(), bson.M{"is_published": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

// GetBlogByID retrieves a blog by its ID
func (br *blogMongoRepo) GetBlogByID(blogID string) (models.Blog, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return models.Blog{}, err
	}

	var blog models.Blog
	err = br.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&blog)
	if err != nil {
		return models.Blog{}, err
	}

	return blog, nil
}

// UpdateBlog updates an existing blog post
func (br *blogMongoRepo) UpdateBlog(blog models.Blog) (models.Blog, error) {
	blog.UpdatedAt = time.Now()

	filter := bson.M{"_id": blog.ID}
	update := bson.M{"$set": blog}

	_, err := br.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return models.Blog{}, err
	}

	return blog, nil
}

// DeleteBlog deletes a blog post
func (br *blogMongoRepo) DeleteBlog(blogID string) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	_, err = br.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	return err
}

// SearchBlogs searches for blogs by title or content
func (br *blogMongoRepo) SearchBlogs(query string) ([]models.Blog, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"is_published": true},
			{"$or": []bson.M{
				{"title": bson.M{"$regex": query, "$options": "i"}},
				{"content": bson.M{"$regex": query, "$options": "i"}},
			}},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := br.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

// FilterBlogs filters blogs by tags, date range, and sort order
func (br *blogMongoRepo) FilterBlogs(tags []string, dateRange [2]string, sortBy string) ([]models.Blog, error) {
	filter := bson.M{"is_published": true}

	// Add tags filter if provided
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}

	// Add date range filter if provided
	if dateRange[0] != "" && dateRange[1] != "" {
		startDate, err := time.Parse("2006-01-02", dateRange[0])
		if err == nil {
			endDate, err := time.Parse("2006-01-02", dateRange[1])
			if err == nil {
				endDate = endDate.Add(24 * time.Hour) // Include the entire end date
				filter["created_at"] = bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				}
			}
		}
	}

	// Set sort order
	var sortField string
	var sortOrder int
	switch sortBy {
	case "title":
		sortField = "title"
		sortOrder = 1
	case "created_at":
		sortField = "created_at"
		sortOrder = -1
	case "updated_at":
		sortField = "updated_at"
		sortOrder = -1
	case "view_count":
		sortField = "view_count"
		sortOrder = -1
	case "likes":
		sortField = "likes"
		sortOrder = -1
	default:
		sortField = "created_at"
		sortOrder = -1
	}

	opts := options.Find().SetSort(bson.D{{Key: sortField, Value: sortOrder}})
	cursor, err := br.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err = cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}

// IncrementViewCount increments the view count of a blog
func (br *blogMongoRepo) IncrementViewCount(blogID string) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"view_count": 1}}

	_, err = br.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// UpdateLikes updates the likes count of a blog
func (br *blogMongoRepo) UpdateLikes(blogID string, increment bool) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	var update bson.M
	if increment {
		update = bson.M{"$inc": bson.M{"likes": 1}}
	} else {
		update = bson.M{"$inc": bson.M{"likes": -1}}
	}

	_, err = br.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// UpdateDislikes updates the dislikes count of a blog
func (br *blogMongoRepo) UpdateDislikes(blogID string, increment bool) error {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	var update bson.M
	if increment {
		update = bson.M{"$inc": bson.M{"dislikes": 1}}
	} else {
		update = bson.M{"$inc": bson.M{"dislikes": -1}}
	}

	_, err = br.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// AddComment adds a comment to a blog post
func (br *blogMongoRepo) AddComment(blogID string, comment models.Comment) (models.Comment, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return models.Comment{}, err
	}

	comment.ID = primitive.NewObjectID()
	comment.BlogID = blogID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	filter := bson.M{"_id": objectID}
	update := bson.M{"$push": bson.M{"comments": comment}}

	_, err = br.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

// GetComments retrieves all comments for a blog post
func (br *blogMongoRepo) GetComments(blogID string) ([]models.Comment, error) {
	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, err
	}

	var blog models.Blog
	err = br.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&blog)
	if err != nil {
		return nil, err
	}

	return blog.Comments, nil
}
