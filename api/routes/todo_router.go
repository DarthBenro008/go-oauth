package routes

import (
	"fmt"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			_ = c.JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
		}
		requestBody.SetCreatedAt()
		requestBody.SetUpdatedAt()
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		requestBody.User, _ = primitive.ObjectIDFromHex(id)
		result, dberr := service.InsertTodo(&requestBody)
		return c.JSON(&fiber.Map{
			"status": result,
			"error":  dberr,
		})
	}
}

func updateTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Todo
		err := c.BodyParser(&requestBody)
		if err != nil {
			_ = c.JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
		}
		requestBody.SetUpdatedAt()
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		requestBody.User, _ = primitive.ObjectIDFromHex(id)
		result, dberr := service.UpdateTodo(&requestBody)
		return c.JSON(&fiber.Map{
			"status": result,
			"error":  dberr,
		})
	}
}

func removeTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.DeleteRequest
		err := c.BodyParser(&requestBody)
		TodoID := requestBody.ID
		if err != nil {
			_ = c.JSON(&fiber.Map{
				"status": false,
				"error":  err,
			})
		}
		dberr := service.RemoveTodo(TodoID)
		if dberr != nil {
			_ = c.JSON(&fiber.Map{
				"status": false,
				"error":  err,
			})
		}
		return c.JSON(&fiber.Map{
			"status":  false,
			"message": "updated successfully",
		})
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
			_ = c.JSON(&fiber.Map{
				"status": false,
				"error":  err,
			})
		}
		return c.JSON(&fiber.Map{
			"status": true,
			"Todos":  fetched,
		})
	}
}
