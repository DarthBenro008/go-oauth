package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"log"
	"oauthserver/api/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	auth := app.Group("/auth")
	routes.AuthRouter(auth)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	_ = app.Listen(":3000")
}
