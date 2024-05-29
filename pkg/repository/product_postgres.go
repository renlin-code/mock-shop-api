package repository

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
)

type ProductPostgres struct {
	db *sqlx.DB
	s  *storage.Storage
}

func newProductPostgres(db *sqlx.DB, s *storage.Storage) *ProductPostgres {
	return &ProductPostgres{db, s}
}

func (r *ProductPostgres) GetAll(limit, offset int, search string) ([]domain.Product, error) {
	var products []domain.Product

	query := fmt.Sprintf("SELECT * FROM %s WHERE available=true", productsTable)
	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", search)
	}
	query += " ORDER BY id LIMIT $1 OFFSET $2"

	err := r.db.Select(&products, query, limit, offset)
	if err == sql.ErrNoRows {
		return products, errors_handler.NoRows()
	}

	return products, err
}

func (r *ProductPostgres) GetById(id int) (domain.Product, error) {
	var product domain.Product

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 AND available=true", productsTable)

	err := r.db.Get(&product, query, id)
	if err == sql.ErrNoRows {
		return product, errors_handler.NoRows()
	}

	return product, err
}

func (r *ProductPostgres) CreateProduct(input domain.CreateProductInput, file multipart.File) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	query := fmt.Sprintf(`INSERT INTO %s (
		category_id,
		name,
		description,
		image_url,
		available,
		price,
		undiscounted_price,
		stock
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, productsTable)

	row := tx.QueryRow(query, input.CategoryId, input.Name, input.Description, "", input.Available, input.Price, input.UndiscountedPrice, input.Stock)
	if err := row.Scan(&id); err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "foreign_key_violation" {
			return 0, errors_handler.ForeignKeyViolation()
		}
		return 0, err
	}

	var url string
	if input.ImgFile != nil {
		url, err = r.s.UploadProductImage(id, input.ImgFile, file)
		if err != nil {
			return 0, err
		}
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET image_url=$1 WHERE id=$2", productsTable)
	_, err = tx.Exec(updateQuery, url, id)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (r *ProductPostgres) UpdateProduct(productId int, input domain.UpdateProductInput, file multipart.File) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.CategoryId != nil {
		setValues = append(setValues, fmt.Sprintf("category_id=$%d", argId))
		args = append(args, *input.CategoryId)
		argId++
	}

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Stock != nil {
		setValues = append(setValues, fmt.Sprintf("stock=$%d", argId))
		args = append(args, *input.Stock)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.UndiscountedPrice != nil {
		setValues = append(setValues, fmt.Sprintf("undiscounted_price=$%d", argId))
		args = append(args, *input.UndiscountedPrice)
		argId++
	}

	if input.Available != nil {
		setValues = append(setValues, fmt.Sprintf("available=$%d", argId))
		args = append(args, *input.Available)
		argId++
	}

	if input.ImgFile != nil && file != nil {
		url, err := r.s.UploadProductImage(productId, input.ImgFile, file)

		if err != nil {
			return err
		}

		setValues = append(setValues, fmt.Sprintf("image_url=$%d", argId))
		args = append(args, url)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING id", productsTable, setQuery, argId)
	args = append(args, productId)

	var id int
	err = tx.QueryRow(query, args...).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors_handler.NoRows()
		}

		return err
	}

	return tx.Commit()
}
