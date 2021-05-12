package routes

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"oauthserver/api/presenter"
	"oauthserver/pkg/user"
)

func UserRouter(app fiber.Router, service user.Service) {
	app.Get("/user", getUser(service))
}

func getUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		userData, err := service.GetUser(id)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenter.Failure("User not found"))
		}
		return c.JSON(presenter.Success(userData, "User data found"))
	}
}
