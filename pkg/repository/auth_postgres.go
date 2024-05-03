package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
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
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}
