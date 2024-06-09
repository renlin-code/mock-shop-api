package repository

import (
	"database/sql/driver"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newCategoryPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		limit  int
		offset int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []domain.Category
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "available", "image_url"}).
					AddRow(1, "category name 1", "category description 1", true, "https://test.back.com/data/categories/1/img1.png").
					AddRow(2, "category name 2", "category description 2", true, "https://test.back.com/data/categories/2/img2.png").
					AddRow(3, "category name 3", "category description 3", true, "https://test.back.com/data/categories/3/img3.png")
				mock.ExpectQuery("SELECT (.+) FROM categories").
					WithArgs(3, 0).WillReturnRows(rows)
			},
			input: args{3, 0},
			want: []domain.Category{
				{Id: 1, Name: "category name 1", Description: "category description 1", Available: true, ImageUrl: "https://test.back.com/data/categories/1/img1.png"},
				{Id: 2, Name: "category name 2", Description: "category description 2", Available: true, ImageUrl: "https://test.back.com/data/categories/2/img2.png"},
				{Id: 3, Name: "category name 3", Description: "category description 3", Available: true, ImageUrl: "https://test.back.com/data/categories/3/img3.png"},
			},
		},
		{
			name: "Ok without pagination",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "available", "image_url"}).
					AddRow(1, "category name 1", "category description 1", true, "https://test.back.com/data/categories/1/img1.png").
					AddRow(2, "category name 2", "category description 2", true, "https://test.back.com/data/categories/2/img2.png").
					AddRow(3, "category name 3", "category description 3", true, "https://test.back.com/data/categories/3/img3.png")
				mock.ExpectQuery("SELECT (.+) FROM categories").
					WithArgs(0, 0).WillReturnRows(rows)
			},
			input: args{0, 0},
			want: []domain.Category{
				{Id: 1, Name: "category name 1", Description: "category description 1", Available: true, ImageUrl: "https://test.back.com/data/categories/1/img1.png"},
				{Id: 2, Name: "category name 2", Description: "category description 2", Available: true, ImageUrl: "https://test.back.com/data/categories/2/img2.png"},
				{Id: 3, Name: "category name 3", Description: "category description 3", Available: true, ImageUrl: "https://test.back.com/data/categories/3/img3.png"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll(tt.input.limit, tt.input.offset, "")
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

func TestGetCategoryById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newCategoryPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		id int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    domain.Category
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "available", "image_url"}).
					AddRow(1, "Category name", "Category description", true, "https://test.back.com/data/categories/1/img1.png")
				mock.ExpectQuery("SELECT (.+) FROM categories").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{1},
			want: domain.Category{
				Id:          1,
				Name:        "Category name",
				Description: "Category description",
				Available:   true,
				ImageUrl:    "https://test.back.com/data/categories/1/img1.png",
			},
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "available", "image_url"})
				mock.ExpectQuery("SELECT (.+) FROM categories").
					WithArgs(0).WillReturnRows(rows)
			},
			input:   args{0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetById(tt.input.id)
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

func TestGetCategoryProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newCategoryPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		categoryId int
		limit      int
		offset     int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []domain.Product
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "category_id", "name", "description", "price", "undiscounted_price", "stock", "available", "image_url"}).
					AddRow(1, 1, "product name 1", "product description 1", 0.99, 1.29, 12, true, "https://test.back.com/data/products/1/img1.png").
					AddRow(2, 1, "product name 2", "product description 2", 1.99, 1.99, 2, true, "https://test.back.com/data/products/2/img1.png").
					AddRow(3, 1, "product name 3", "product description 3", 108.49, 126.99, 27, true, "https://test.back.com/data/products/3/img1.png")
				mock.ExpectQuery("SELECT (.+) FROM products").
					WithArgs(1, 3, 0).WillReturnRows(rows)
			},
			input: args{1, 3, 0},
			want: []domain.Product{
				{Id: 1, CategoryId: 1, Name: "product name 1", Description: "product description 1", Price: 0.99, UndiscountedPrice: 1.29, Stock: 12, Available: true, ImageUrl: "https://test.back.com/data/products/1/img1.png"},
				{Id: 2, CategoryId: 1, Name: "product name 2", Description: "product description 2", Price: 1.99, UndiscountedPrice: 1.99, Stock: 2, Available: true, ImageUrl: "https://test.back.com/data/products/2/img1.png"},
				{Id: 3, CategoryId: 1, Name: "product name 3", Description: "product description 3", Price: 108.49, UndiscountedPrice: 126.99, Stock: 27, Available: true, ImageUrl: "https://test.back.com/data/products/3/img1.png"},
			},
		},
		{
			name: "Ok without pagination",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "category_id", "name", "description", "price", "undiscounted_price", "stock", "available", "image_url"}).
					AddRow(1, 1, "product name 1", "product description 1", 0.99, 1.29, 12, true, "https://test.back.com/data/products/1/img1.png").
					AddRow(2, 1, "product name 2", "product description 2", 1.99, 1.99, 2, true, "https://test.back.com/data/products/2/img1.png").
					AddRow(3, 1, "product name 3", "product description 3", 108.49, 126.99, 27, true, "https://test.back.com/data/products/3/img1.png")
				mock.ExpectQuery("SELECT (.+) FROM products").
					WithArgs(1, 0, 0).WillReturnRows(rows)
			},
			input: args{1, 0, 0},
			want: []domain.Product{
				{Id: 1, CategoryId: 1, Name: "product name 1", Description: "product description 1", Price: 0.99, UndiscountedPrice: 1.29, Stock: 12, Available: true, ImageUrl: "https://test.back.com/data/products/1/img1.png"},
				{Id: 2, CategoryId: 1, Name: "product name 2", Description: "product description 2", Price: 1.99, UndiscountedPrice: 1.99, Stock: 2, Available: true, ImageUrl: "https://test.back.com/data/products/2/img1.png"},
				{Id: 3, CategoryId: 1, Name: "product name 3", Description: "product description 3", Price: 108.49, UndiscountedPrice: 126.99, Stock: 27, Available: true, ImageUrl: "https://test.back.com/data/products/3/img1.png"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetProducts(tt.input.categoryId, tt.input.limit, tt.input.offset, "")
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

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newCategoryPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		input domain.CreateCategoryInput
		file  multipart.File
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
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO categories").
					WithArgs("Category name", "Category description", "", true).WillReturnRows(rows)

				mock.ExpectExec("UPDATE categories").
					WithArgs("https://test.back.com/data/categories/1/img1.png", 1).WillReturnResult(driver.ResultNoRows)

				mock.ExpectCommit()
			},
			input: args{
				file: testFile,
				input: domain.CreateCategoryInput{
					Name:        "Category name",
					Description: "Category description",
					Available:   true,
					ImgFile:     testFileHeader,
				},
			},
			want: 1,
		},
		{
			name: "CategoryExists",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO categories").WithArgs("Category name", "Category description", "", true).WillReturnError(errors.New("category already exists"))
				mock.ExpectRollback()
			},
			input: args{
				file: testFile,
				input: domain.CreateCategoryInput{
					Name:        "Category name",
					Description: "Category description",
					Available:   true,
					ImgFile:     testFileHeader,
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateCategory(tt.input.input, tt.input.file)
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

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newCategoryPostgres(sqlx.NewDb(db, "sqlmock"), s)
	type args struct {
		categoryId int
		input      domain.UpdateCategoryInput
		file       multipart.File
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
				mock.ExpectQuery("UPDATE categories SET (.+)").
					WithArgs("new name", "new description", true, "https://test.back.com/data/categories/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				categoryId: 1,
				input: domain.UpdateCategoryInput{
					Name:        stringPointer("new name"),
					Description: stringPointer("new description"),
					Available:   boolPointer(true),
					ImgFile:     testFileHeader,
				},
				file: testFile,
			},
			want: 1,
		},
		{
			name: "OK_NoAvailable",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("UPDATE categories SET (.+)").
					WithArgs("new name", "new description", "https://test.back.com/data/categories/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				categoryId: 1,
				input: domain.UpdateCategoryInput{
					Name:        stringPointer("new name"),
					Description: stringPointer("new description"),
					ImgFile:     testFileHeader,
				},
				file: testFile,
			},
			want: 1,
		},
		{
			name: "OK_NoAvailableAndDescription",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("UPDATE categories SET (.+)").
					WithArgs("new name", "https://test.back.com/data/categories/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				categoryId: 1,
				input: domain.UpdateCategoryInput{
					Name:    stringPointer("new name"),
					ImgFile: testFileHeader,
				},
				file: testFile,
			},
			want: 1,
		},
		{
			name: "OK_NoFileAndDescription",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("UPDATE categories SET (.+)").
					WithArgs("new name", true, 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				categoryId: 1,
				input: domain.UpdateCategoryInput{
					Name:      stringPointer("new name"),
					Available: boolPointer(true),
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateCategory(tt.input.categoryId, tt.input.input, tt.input.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
