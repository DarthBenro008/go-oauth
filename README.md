# Go Oauth Server

This is sample OAuth integration written in GoLang that also uses MongoDB. This is a sample TODO Application where people can Create and update Tasks.

The application makes use of both **Google Authentication** and **GitHub Authentication**

The documentation can be found here:

[Postman Documentation](https://documenter.getpostman.com/view/12292853/TzRUAmao#intro)

## Tech Stack

- [GoFiber](https://github.com/gofiber/fiber)
- [MongoDB Driver](https://github.com/mongodb/mongo-go-driver)

### Architecture

This Application makes use of `Clean Architecture`

This architecture enables us to add more authentication methods in future and also we can shift to a different database completely as the API makes use of services and the database implementation is written in the repository.

## Usage

- Clone this repo
- Create a .env file like the following sample

```bash
GOOGLE_OAUTH_CLIENT_ID=
GOOGLE_OAUTH_CLIENT_SECRET=
GITHUB_OAUTH_CLIENT_ID=
GITHUB_OAUTH_CLIENT_SECRET=
JWT_SECRET=
MONGO_URI=
```

### GitHub OAuth2

- Follow the documentation given [here](https://docs.github.com/en/developers/apps/authorizing-oauth-apps)

### Google OAuth2

- Follow the documentation given [here](https://developers.google.com/identity/protocols/oauth2/javascript-implicit-flow)

- You can run `docker-compose up` this will spin up the mongoose database also for you.

## Contributions

- Feel Free to Open a PR/Issue for any feature or bug(s).
- Make sure you follow the community guidelines!
- Feel free to open an issue to ask a question/discuss anything with regards to any improvements in the implementation.
- Have a feature request? Open an Issue!

## License

Copyright 2021 Hemanth Krishna

Licensed under MIT License : https://opensource.org/licenses/MIT
