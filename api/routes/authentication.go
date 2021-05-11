package routes

import (
	"github.com/gofiber/fiber/v2"
	"oauthserver/api/service"
)

func AuthRouter(app fiber.Router) {
	app.Get("/google/login", service.GoogleLogin())
	app.Get("/google/callback", service.GoogleCallback())
}
