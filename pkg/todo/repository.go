package todo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauthserver/pkg/entities"
)

type Repository interface {
	CreateTodo(todo *entities.Todo) (*entities.Todo, error)
	ReadTodo(userId string) (*[]entities.Todo, error)
	UpdateTodo(todo *entities.Todo) (*entities.Todo, error)
	DeleteTodo(ID string) error
}
type repository struct {
	Collection *mongo.Collection
}

func (r repository) CreateTodo(todo *entities.Todo) (*entities.Todo, error) {
	todo.ID = primitive.NewObjectID()
	_, err := r.Collection.InsertOne(context.Background(), todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r repository) ReadTodo(userId string) (*[]entities.Todo, error) {
	var Todos []entities.Todo
	cursor, err := r.Collection.Find(context.Background(), bson.M{
		"user": userId,
	})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var todo entities.Todo
		_ = cursor.Decode(&todo)
		Todos = append(Todos, todo)
	}
	return &Todos, nil
}

func (r repository) UpdateTodo(todo *entities.Todo) (*entities.Todo, error) {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": todo.ID}, bson.M{"$set": todo})
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r repository) DeleteTodo(ID string) error {
	TodoID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(context.Background(), bson.M{"_id": TodoID})
	if err != nil {
		return err
	}
	return nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
