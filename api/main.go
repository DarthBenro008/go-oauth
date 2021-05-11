package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"oauthserver/api/presenter"
	"oauthserver/api/routes"
	"oauthserver/pkg/todo"
	"oauthserver/pkg/user"
	"os"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := DatabaseConnection()
	if err != nil {
		log.Fatal("Database Connection Error $s", err)
	}
	fmt.Println("Database connection success!")

	todoCollection := db.Collection("todos")
	todoRepo := todo.NewRepo(todoCollection)
	todoService := todo.NewService(todoRepo)

	userCollection := db.Collection("users")
	userRepo := user.NewRepo(userCollection)
	userService := user.NewService(userRepo)


	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(presenter.Success("The Litmus Task API is running", "Author: Hemanth Krishna [@DarthBenro008]"))
	})
	auth := app.Group("/auth")
	routes.AuthRouter(auth, userService)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	todos := app.Group("/api")
	routes.TodoRouter(todos, todoService)

	_ = app.Listen(":3000")
}

func DatabaseConnection() (*mongo.Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/litmusTask"))
	if err != nil {
		return nil, err
	}
	db := client.Database("litmusTask")
	return db, nil
}
