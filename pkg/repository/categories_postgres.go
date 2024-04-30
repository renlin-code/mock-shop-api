package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func newCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db}
}

func (r *CategoryPostgres) GetAll() ([]domain.Category, error) {
	var categories []domain.Category

	query := fmt.Sprintf("SELECT * FROM %s", categoriesTables)

	err := r.db.Select(&categories, query)
	return categories, err
}

func (r *CategoryPostgres) GetById(id int) (domain.Category, error) {
	var category domain.Category

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", categoriesTables)

	err := r.db.Get(&category, query, id)

	return category, err
}
