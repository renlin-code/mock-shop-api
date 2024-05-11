package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func newProductPostgres(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db}
}

func (r *ProductPostgres) GetAll() ([]domain.Product, error) {
	var products []domain.Product

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", productsTable)

	err := r.db.Select(&products, query)
	if err == sql.ErrNoRows {
		return products, errors_handler.NoRows()
	}

	return products, err
}

func (r *ProductPostgres) GetById(id int) (domain.Product, error) {
	var product domain.Product

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", productsTable)

	err := r.db.Get(&product, query, id)
	if err == sql.ErrNoRows {
		return product, errors_handler.NoRows()
	}

	return product, err
}
