package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauthserver/pkg/entities"
)

type Repository interface {
	LoginUser(user *entities.User) (*entities.User, error)
}
type repository struct {
	Collection *mongo.Collection
}

func (r repository) LoginUser(user *entities.User) (*entities.User, error) {
	user.ID = primitive.NewObjectID()
	var result = entities.User{}
	findOneErr := r.Collection.FindOne(context.TODO(), bson.M{
		"email":  user.Email,
		"social": user.Social,
	}).Decode(&result)

	if findOneErr != nil {
		_, err := r.Collection.InsertOne(context.Background(), user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return &result, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
