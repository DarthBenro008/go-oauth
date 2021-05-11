package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
	"oauthserver/api/presenter"
	"oauthserver/pkg/entities"
	"oauthserver/pkg/user"
	"oauthserver/pkg/utils"
	"os"
)

var (
	gitHubState = "randomState"
)

type GitHubResponseStruct struct {
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type GithubEmailResponse []struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func oAuthGitHubConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/github/callback",
		ClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
}

func GitHubLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tempState, err := utils.GenerateRandomString()
		gitHubState = tempState
		if err != nil {
			fmt.Println(err)
			return c.SendString("Some error has occurred.")
		}
		url := oAuthGitHubConfig().AuthCodeURL(gitHubState, oauth2.AccessTypeOnline)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	}
}

func GithubCallback(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.FormValue("state") != gitHubState {
			return c.Redirect("/", http.StatusTemporaryRedirect)
		}

		token, err := oAuthGitHubConfig().Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			fmt.Print(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}
		req.Header.Set("Authorization", "token "+token.AccessToken)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return  c.JSON(presenter.Failure(err))
		}
		githubResponse := GitHubResponseStruct{}
		err = json.NewDecoder(resp.Body).Decode(&githubResponse)
		if err != nil {
			fmt.Println(err)
		}
		if githubResponse.Email == "" {
			gitHubEmails := GithubEmailResponse{}
			req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return  c.JSON(presenter.Failure(err))
			}
			req.Header.Set("Authorization", "token "+token.AccessToken)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				c.Status(http.StatusInternalServerError)
				return  c.JSON(presenter.Failure(err))
			}
			err = json.NewDecoder(resp.Body).Decode(&gitHubEmails)
			if err != nil {
				fmt.Println(err)
			}
			githubResponse.Email = gitHubEmails[0].Email
		}
		var userData = entities.User{
			Name:        githubResponse.Name,
			Email:       githubResponse.Email,
			AccessToken: token.AccessToken,
			Social:      "github",
			Picture:     githubResponse.AvatarURL,
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
		return c.JSON(presenter.TokenResponse(jwtToken))
	}

}
