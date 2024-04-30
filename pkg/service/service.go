package service

import (
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type Authorization interface {
	UserSignUp(name, email string) error
	CreateUser(token, password string) (int, error)
	GenerateAuthToken(email, password string) (string, error)
	ParseAuthToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos),
	}
}
