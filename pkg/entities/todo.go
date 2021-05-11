package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Todo struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Done        bool               `json:"done"`
	CreatedAt   int64              `json:"createdAt"`
	UpdatedAt   int64              `json:"updatedAt"`
	User        primitive.ObjectID `json:"user"`
}

func (todo *Todo) SetCreatedAt() {
	todo.CreatedAt = time.Now().Unix()
}

func (todo *Todo) SetUpdatedAt() {
	todo.UpdatedAt = time.Now().Unix()
}

type DeleteRequest struct {
	ID string `json:"id" binding:"required,gte=1"`
}
