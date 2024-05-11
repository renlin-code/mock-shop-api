package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func newAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, profile_image) VALUES ($1, $2, $3, '') RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "unique_violation" {
			return 0, errors_handler.AlreadyExists("user")
		}
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	if err == sql.ErrNoRows {
		return user, errors_handler.NoRows()
	}
	return user, err
}

func (r *AuthPostgres) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, name, email FROM %s WHERE email=$1", usersTable)
	err := r.db.Get(&user, query, email)
	if err == sql.ErrNoRows {
		return user, errors_handler.NoRows()
	}
	return user, err
}

func (r *AuthPostgres) UpdatePassword(userId int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE id=$2", usersTable)

	_, err := r.db.Exec(query, password, userId)
	if err == sql.ErrNoRows {
		return errors_handler.NoRows()
	}
	return err
}
