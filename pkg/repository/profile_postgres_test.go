package repository

import (
	"mime/multipart"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		BaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProfilePostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		userId int
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
				mock.ExpectQuery("SELECT id, name, email, profile_image FROM users").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{1},
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
				mock.ExpectQuery("SELECT id, name, email, profile_image FROM users").
					WithArgs(0).WillReturnRows(rows)
			},
			input:   args{0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetProfile(tt.input.userId)
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

func TestUpdateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		BaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProfilePostgres(sqlx.NewDb(db, "sqlmock"), s)
	type args struct {
		userId int
		input  domain.UpdateProfileInput
		file   multipart.File
	}

	file, err := os.CreateTemp("", "img1.png")
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	testFileHeader := &multipart.FileHeader{
		Filename: "img1.png",
	}

	testFile, err := os.Open(file.Name())
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}
	defer testFile.Close()

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("UPDATE users SET (.+)").
					WithArgs("new name", "https://test.back.com/data/users/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				userId: 1,
				input: domain.UpdateProfileInput{
					Name:           stringPointer("new name"),
					ProfileImgFile: testFileHeader,
				},
				file: testFile,
			},
			want: 1,
		},
		{
			name: "OK_NoFile",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("UPDATE users SET (.+)").
					WithArgs("new name", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				userId: 1,
				input: domain.UpdateProfileInput{
					Name: stringPointer("new name"),
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateProfile(tt.input.userId, tt.input.input, tt.input.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		BaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProfilePostgres(sqlx.NewDb(db, "sqlmock"), s)

	userId := 1
	products := []domain.CreateOrderInputProduct{
		{Id: 1, Quantity: 5},
		{Id: 2, Quantity: 3},
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(123)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(userId, sqlmock.AnyArg()).
		WillReturnRows(rows)
	mock.ExpectPrepare("INSERT INTO ordered_products")

	for _, product := range products {
		rows = sqlmock.NewRows([]string{"name", "description", "price", "undiscounted_price", "image_url", "stock"}).
			AddRow("Product1", "Description1", 10.5, 12.0, "image1.jpg", 20)

		mock.ExpectQuery("UPDATE products SET stock").
			WithArgs(product.Quantity, product.Id).
			WillReturnRows(rows)

		mock.ExpectExec("INSERT INTO ordered_products").
			WithArgs(123, product.Id, "Product1", "Description1", 10.5, 12.0, "image1.jpg", product.Quantity).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	orderId, err := r.CreateOrder(userId, products)
	assert.NoError(t, err)
	assert.Equal(t, 123, orderId)

	assert.NoError(t, mock.ExpectationsWereMet())
}
