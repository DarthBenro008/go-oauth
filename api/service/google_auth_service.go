package service

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"oauthserver/pkg/utils"
	"os"
)

var (
	state = "holderState"
)

func oAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tempState, err := utils.GenerateRandomString()
		state = tempState
		if err != nil {
			fmt.Println(err)
			return c.SendString("Some error has occurred.")
		}
		url := oAuthConfig().AuthCodeURL(state)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	}
}

func GoogleCallback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.FormValue("state") != state {
			return c.Redirect("/", http.StatusTemporaryRedirect)
		}
		token, err := oAuthConfig().Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			fmt.Print(err)
			return c.SendString("Sorry, something went wrong ")
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			return c.SendString("Cannot get your details bro")
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		return c.Send(content)
	}
}
