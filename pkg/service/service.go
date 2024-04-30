package service

import (
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type Authorization interface {
	UserSignUp(name, email string) error
	CreateUser(token, password string) (int, error)
	GenerateAuthToken(email, password string) (string, error)
	ParseAuthToken(token string) (int, error)
}

type Category interface {
	GetAll() ([]domain.Category, error)
	GetById(id int) (domain.Category, error)
}

type Service struct {
	Authorization
	Category
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos),
		Category:      newCategoryService(repos),
	}
}
