package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func newCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db}
}

func (r *CategoryPostgres) GetAll() ([]domain.Category, error) {
	var categories []domain.Category

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", categoriesTables)

	err := r.db.Select(&categories, query)
	if err == sql.ErrNoRows {
		return categories, errors_handler.NoRows()
	}

	return categories, err
}

func (r *CategoryPostgres) GetById(id int) (domain.Category, error) {
	var category domain.Category

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", categoriesTables)

	err := r.db.Get(&category, query, id)
	if err == sql.ErrNoRows {
		return category, errors_handler.NoRows()
	}

	return category, err
}

func (r *CategoryPostgres) GetProducts(categoryId int) ([]domain.Product, error) {
	var products []domain.Product

	query := fmt.Sprintf("SELECT * FROM %s WHERE category_id=$1 ORDER BY id", productsTable)

	err := r.db.Select(&products, query, categoryId)
	if err == sql.ErrNoRows {
		return products, errors_handler.NoRows()
	}

	return products, err
}
