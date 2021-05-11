package todo

import "oauthserver/pkg/entities"

type Service interface {
	InsertTodo(todo *entities.Todo) (*entities.Todo, error)
	FetchTodos(userId string) (*[]entities.Todo, error)
	UpdateTodo(todo *entities.Todo) (*entities.Todo, error)
	RemoveTodo(ID string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertTodo(todo *entities.Todo) (*entities.Todo, error) {
	return s.repository.CreateTodo(todo)
}
func (s *service) FetchTodos(userId string) (*[]entities.Todo, error) {
	return s.repository.ReadTodo(userId)

}
func (s *service) UpdateTodo(todo *entities.Todo) (*entities.Todo, error) {
	return s.repository.UpdateTodo(todo)
}
func (s *service) RemoveTodo(ID string) error {
	return s.repository.DeleteTodo(ID)
}
