package entities

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	AccessToken string             `json:"accessToken"`
	Social      string             `json:"social"`
	Picture     string             `json:"picture"`
}

func (user *User) GetSignedJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID.Hex()
	claims["exp"] = time.Now().Add(time.Minute * 300).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}
