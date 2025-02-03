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
		for _, endpoint := range controllers.PublicEndpoints {
			public.Handle(endpoint.Method, endpoint.Path, endpoint.HandlerFunction(db))
		}
	}

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(db))
	{
		for _, endpoint := range controllers.PrivateEndpoints {
			auth.Handle(endpoint.Method, endpoint.Path, endpoint.HandlerFunction(db))
		}
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
