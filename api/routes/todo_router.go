package routes

import (
	"fmt"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"oauthserver/api/presenter"
	"oauthserver/pkg/entities"
	"oauthserver/pkg/todo"
)

func TodoRouter(app fiber.Router, service todo.Service) {
	app.Get("/todos", getTodos(service))
	app.Post("/todos", addTodo(service))
	app.Put("/todos", updateTodo(service))
	app.Delete("/todos", removeTodo(service))
}

func addTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Todo
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}
		requestBody.SetCreatedAt()
		requestBody.SetUpdatedAt()
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		requestBody.User, _ = primitive.ObjectIDFromHex(id)
		result, dberr := service.InsertTodo(&requestBody)
		return c.JSON(presenter.Success(result, dberr))
	}
}

func updateTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Todo
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}
		requestBody.SetUpdatedAt()
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		requestBody.User, _ = primitive.ObjectIDFromHex(id)
		result, dberr := service.UpdateTodo(&requestBody)
		return c.JSON(presenter.Success(result, dberr))
	}
}

func removeTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.DeleteRequest
		err := c.BodyParser(&requestBody)
		TodoID := requestBody.ID
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}
		dberr := service.RemoveTodo(TodoID)
		if dberr != nil {
			_ = c.JSON(presenter.Failure(err))
		}
		return c.JSON(presenter.Success("Removed Successfully", dberr))
	}
}

func getTodos(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		fmt.Println("This is the token received", id)
		fetched, err := service.FetchTodos(id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.Failure(err))
		}
		return c.JSON(presenter.Success(fetched, "Fetched todos successfully"))
	}
}
