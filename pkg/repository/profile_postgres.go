package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type ProfilePostgres struct {
	db *sqlx.DB
}

func newProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{db}
}

func (r *ProfilePostgres) GetProfile(userId int) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, name, email, profile_image FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *ProfilePostgres) UpdateProfile(userId int, input domain.UpdateProfileInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.ProfileImg != nil {
		setValues = append(setValues, fmt.Sprintf("profile_image=$%d", argId))
		args = append(args, *input.ProfileImg)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", usersTable, setQuery, argId)
	args = append(args, userId)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfilePostgres) UpdatePassword(userId int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE id=$2", usersTable)

	fmt.Println(query)
	_, err := r.db.Exec(query, password, userId)
	return err
}
