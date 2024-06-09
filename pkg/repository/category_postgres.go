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

type CategoryPostgres struct {
	db *sqlx.DB
	s  *storage.Storage
}

func newCategoryPostgres(db *sqlx.DB, s *storage.Storage) *CategoryPostgres {
	return &CategoryPostgres{db, s}
}

func (r *CategoryPostgres) GetAll(limit, offset int, search string) ([]domain.Category, error) {
	var categories []domain.Category

	query := fmt.Sprintf("SELECT * FROM %s WHERE available=true", categoriesTables)
	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", search)
	}
	query += " ORDER BY id LIMIT $1 OFFSET $2"

	err := r.db.Select(&categories, query, limit, offset)
	if err == sql.ErrNoRows {
		return categories, errors_handler.NoRows()
	}

	return categories, err
}

func (r *CategoryPostgres) GetById(id int) (domain.Category, error) {
	var category domain.Category

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 AND available=true", categoriesTables)

	err := r.db.Get(&category, query, id)
	if err == sql.ErrNoRows {
		return category, errors_handler.NoRows()
	}

	return category, err
}

func (r *CategoryPostgres) GetFilePath(categoryId int, fileName string) string {
	return r.s.Category.GetFilePath(categoryId, fileName)
}

func (r *CategoryPostgres) GetProducts(categoryId, limit, offset int, search string) ([]domain.Product, error) {
	var products []domain.Product

	query := fmt.Sprintf("SELECT * FROM %s WHERE category_id=$1 AND available=true", productsTable)
	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", search)
	}
	query += " ORDER BY id LIMIT $2 OFFSET $3"

	err := r.db.Select(&products, query, categoryId, limit, offset)
	if err == sql.ErrNoRows {
		return products, errors_handler.NoRows()
	}

	return products, err
}

func (r *CategoryPostgres) CreateCategory(input domain.CreateCategoryInput, file multipart.File) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	query := fmt.Sprintf(`INSERT INTO %s (
		name, 
		description, 
		image_url, 
		available
	) VALUES ($1, $2, $3, $4) RETURNING id`, categoriesTables)

	row := tx.QueryRow(query, input.Name, input.Description, "", input.Available)
	if err := row.Scan(&id); err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "unique_violation" {
			return 0, errors_handler.AlreadyExists("category")
		}
		return 0, err
	}

	var url string
	if input.ImgFile != nil {
		url, err = r.s.UploadCategoryImage(id, input.ImgFile, file)
		if err != nil {
			return 0, err
		}
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET image_url=$1 WHERE id=$2", categoriesTables)
	_, err = tx.Exec(updateQuery, url, id)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (r *CategoryPostgres) UpdateCategory(categoryId int, input domain.UpdateCategoryInput, file multipart.File) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

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

	if input.Available != nil {
		setValues = append(setValues, fmt.Sprintf("available=$%d", argId))
		args = append(args, *input.Available)
		argId++
	}

	if input.ImgFile != nil && file != nil {
		url, err := r.s.UploadCategoryImage(categoryId, input.ImgFile, file)

		if err != nil {
			return err
		}

		setValues = append(setValues, fmt.Sprintf("image_url=$%d", argId))
		args = append(args, url)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING id", categoriesTables, setQuery, argId)
	args = append(args, categoryId)

	var id int
	err = tx.QueryRow(query, args...).Scan(&id)
	pqErr, ok := err.(*pq.Error)
	if ok && pqErr.Code.Name() == "unique_violation" {
		return errors_handler.AlreadyExists("category")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return errors_handler.NoRows()
		}

		return err
	}

	return tx.Commit()
}
