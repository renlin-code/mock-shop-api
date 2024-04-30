package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email, password string) (domain.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
	}
}
