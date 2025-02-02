package startup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	middleware "main/authentication"
	"main/constants"
	"main/controllers"
	"main/models"
	"os"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	public := router.Group("/")
	{
		public.POST(controllers.UserSignUpEndpoint.Path, controllers.UserSignUpEndpoint.HandlerFunction(db))
		public.POST(controllers.UserLoginEndpoint.Path, controllers.UserLoginEndpoint.HandlerFunction(db))

		public.GET(controllers.ViewUserProfileEndpoint.Path, controllers.ViewUserProfileEndpoint.HandlerFunction(db))

		public.GET(controllers.SearchUserEndpoint.Path, controllers.SearchUserEndpoint.HandlerFunction(db))
		public.GET(controllers.SearchPostEndpoint.Path, controllers.SearchPostEndpoint.HandlerFunction(db))
	}

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(db))
	{
		auth.POST(controllers.FollowUserEndpoint.Path, controllers.FollowUserEndpoint.HandlerFunction(db))
		auth.DELETE(controllers.UnfollowUserEndpoint.Path, controllers.UnfollowUserEndpoint.HandlerFunction(db))
		auth.GET(controllers.GetFollowersProfileEndpoint.Path, controllers.GetFollowersProfileEndpoint.HandlerFunction(db))
		auth.GET(controllers.GetFollowingProfileEndpoint.Path, controllers.GetFollowingProfileEndpoint.HandlerFunction(db))

		auth.PUT(controllers.EditUserProfileEndpoint.Path, controllers.EditUserProfileEndpoint.HandlerFunction(db))

		auth.POST(controllers.CreatePostEndpoint.Path, controllers.CreatePostEndpoint.HandlerFunction(db))
		auth.PUT(controllers.EditPostEndpoint.Path, controllers.EditPostEndpoint.HandlerFunction(db))
		auth.DELETE(controllers.DeletePostEndpoint.Path, controllers.DeletePostEndpoint.HandlerFunction(db))
		auth.GET(controllers.GetSpecificPostEndpoint.Path, controllers.GetSpecificPostEndpoint.HandlerFunction(db))
		auth.GET(controllers.GetAllPostsByUserIDEndpoint.Path, controllers.GetAllPostsByUserIDEndpoint.HandlerFunction(db))
		auth.GET(controllers.GetAllPostsEndpoint.Path, controllers.GetAllPostsEndpoint.HandlerFunction(db))
	}

	return router
}

func StartRoutes(db *gorm.DB) error {
	return SetupRouter(db).Run()
}

func StartDatabase() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil
	}

	host := os.Getenv("HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")

	envVariables := []string{host, user, password, dbname, port}

	for _, envVar := range envVariables {
		if envVar == constants.EMPTY {
			log.Fatal("One or more database environment variables are not set")
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatalf("failed to enable uuid-ossp extension: %v", err)
	}

	migrateSchemas(db)

	return db
}

func migrateSchemas(db *gorm.DB) {
	err := db.AutoMigrate(&models.Post{}, &models.Follow{}, &models.Like{}, &models.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
