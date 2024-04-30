package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func newProductPostgres(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db}
}

func (r *ProductPostgres) GetAll() ([]domain.Product, error) {
	var products []domain.Product

	query := fmt.Sprintf("SELECT * FROM %s", productsTable)

	err := r.db.Select(&products, query)
	return products, err
}

func (r *ProductPostgres) GetById(id int) (domain.Product, error) {
	var product domain.Product

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", productsTable)

	err := r.db.Get(&product, query, id)

	return product, err
}
