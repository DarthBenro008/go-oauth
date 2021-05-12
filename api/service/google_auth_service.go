package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"oauthserver/api/presenter"
	"oauthserver/pkg/entities"
	"oauthserver/pkg/user"
	"oauthserver/pkg/utils"
	"os"
)

var (
	state = "holderState"
)

type googleAuthResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func oAuthGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
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
		url := oAuthGoogleConfig().AuthCodeURL(state)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	}
}

func GoogleCallback(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.FormValue("state") != state {
			return c.Redirect("/", http.StatusTemporaryRedirect)
		}
		token, err := oAuthGoogleConfig().Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			fmt.Print(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			return c.SendString("Cannot get your details bro")
		}
		defer resp.Body.Close()
		googleResponse := googleAuthResponse{}
		err = json.NewDecoder(resp.Body).Decode(&googleResponse)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}
		var userData = entities.User{
			Name:        googleResponse.GivenName + " " + googleResponse.FamilyName,
			Email:       googleResponse.Email,
			AccessToken: token.AccessToken,
			Social:      "google",
			Picture:     googleResponse.Picture,
		}
		signedInUser, err := userService.LoginUser(&userData)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}
		jwtToken, err := signedInUser.GetSignedJWT()
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}
		userData.AccessToken = ""
		return c.JSON(presenter.TokenResponse(jwtToken, userData))
	}
}
