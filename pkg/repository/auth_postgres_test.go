package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	r := newAuthPostgres(sqlx.NewDb(db, "sqlmock"))

	tests := []struct {
		name    string
		mock    func()
		input   domain.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("Alice", "alice@example.com", "password").WillReturnRows(rows)
			},
			input: domain.User{
				Name:       "Alice",
				Email:      "alice@example.com",
				Password:   "password",
				ProfileImg: "",
			},
			want: 1,
		},
		{
			name: "Empty Field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("Alice", "alice@example.com", "").WillReturnRows(rows)
			},
			input: domain.User{
				Name:       "Alice",
				Email:      "alice@example.com",
				Password:   "",
				ProfileImg: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	r := newAuthPostgres(sqlx.NewDb(db, "sqlmock"))

	type args struct {
		email    string
		password string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    domain.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "profile_image"}).
					AddRow(1, "Alice", "alice@example.com", "password", "https://url-image.png")
				mock.ExpectQuery("SELECT id FROM users").
					WithArgs("alice@example.com", "password").WillReturnRows(rows)
			},
			input: args{"alice@example.com", "password"},
			want: domain.User{
				Id:         1,
				Name:       "Alice",
				Email:      "alice@example.com",
				Password:   "password",
				ProfileImg: "https://url-image.png",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password", "profile_image"})
				mock.ExpectQuery("SELECT id FROM users").
					WithArgs("not", "found").WillReturnRows(rows)
			},
			input:   args{"not", "found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUser(tt.input.email, tt.input.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	r := newAuthPostgres(sqlx.NewDb(db, "sqlmock"))

	type args struct {
		email string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    domain.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email"}).
					AddRow(1, "Alice", "alice@example.com")
				mock.ExpectQuery("SELECT id, name, email FROM users").
					WithArgs("alice@example.com").WillReturnRows(rows)
			},
			input: args{"alice@example.com"},
			want: domain.User{
				Id:         1,
				Name:       "Alice",
				Email:      "alice@example.com",
				Password:   "",
				ProfileImg: "",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email"})
				mock.ExpectQuery("SELECT id, name, email FROM users").
					WithArgs("not found").WillReturnRows(rows)
			},
			input:   args{"not found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByEmail(tt.input.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	r := newAuthPostgres(sqlx.NewDb(db, "sqlmock"))

	type args struct {
		password string
		id       int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    error
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("UPDATE users").
					WithArgs("new password", 1).WillReturnRows(rows)
			},
			input: args{"new password", 1},
			want:  nil,
		},
		{
			name: "Empty Field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("UPDATE users").
					WithArgs("", 1).WillReturnRows(rows)
			},
			input: args{"", 1},

			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdatePassword(tt.input.id, tt.input.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, nil)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
