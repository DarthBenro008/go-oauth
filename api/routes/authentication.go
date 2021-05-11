package routes

import (
	"github.com/gofiber/fiber/v2"
	"oauthserver/api/service"
	"oauthserver/pkg/user"
)

func AuthRouter(app fiber.Router, userService user.Service) {
	app.Get("/google/login", service.GoogleLogin())
	app.Get("/google/callback", service.GoogleCallback(userService))
	app.Get("/github/login", service.GitHubLogin())
	app.Get("/github/callback", service.GithubCallback(userService))
}
