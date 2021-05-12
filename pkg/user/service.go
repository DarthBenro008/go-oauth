package user

import "oauthserver/pkg/entities"

type Service interface {
	LoginUser(user *entities.User) (*entities.User, error)
	GetUser(userId string) (*entities.User, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) LoginUser(user *entities.User) (*entities.User, error) {
	return s.repository.LoginUser(user)
}
func (s *service) GetUser(userId string) (*entities.User, error) {
	return s.repository.GetUser(userId)
}
