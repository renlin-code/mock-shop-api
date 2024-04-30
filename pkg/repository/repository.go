package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email, password string) (domain.User, error)
}

type Category interface {
	GetAll() ([]domain.Category, error)
	GetById(id int) (domain.Category, error)
}

type Repository struct {
	Authorization
	Category
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
		Category:      newCategoryPostgres(db),
	}
}
