package main

import (
	"blog-api/Delivery/controllers"
	"blog-api/Delivery/routers"
	"blog-api/Infrastructure/database"
	"blog-api/Infrastructure/repositories"
	usecases "blog-api/Usecases"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	mongo, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}

	// dbName := os.Getenv("DB_NAME")

	blogCollection := mongo.GetCollection("blog_db", "blogs")

	blogRepo := repositories.NewBlogRepository(blogCollection)

	blogUseCase := usecases.NewBlogUseCase(blogRepo)

	blogController := controllers.NewBlogController(blogUseCase)

	r := routers.SetupRoutes(blogController)
	r.Run(":8080")

}
