package repository

import (
	"database/sql/driver"
	"errors"
	"math"
	"mime/multipart"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProductPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		limit  int
		offset int
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
					AddRow(3, 2, "product name 3", "product description 3", 108.49, 126.99, 27, true, "https://test.back.com/data/products/3/img1.png")
				mock.ExpectQuery("SELECT (.+) FROM products").
					WithArgs(3, 0).WillReturnRows(rows)
			},
			input: args{3, 0},
			want: []domain.Product{
				{Id: 1, CategoryId: 1, Name: "product name 1", Description: "product description 1", Price: 0.99, UndiscountedPrice: 1.29, Stock: 12, Available: true, ImageUrl: "https://test.back.com/data/products/1/img1.png"},
				{Id: 2, CategoryId: 1, Name: "product name 2", Description: "product description 2", Price: 1.99, UndiscountedPrice: 1.99, Stock: 2, Available: true, ImageUrl: "https://test.back.com/data/products/2/img1.png"},
				{Id: 3, CategoryId: 2, Name: "product name 3", Description: "product description 3", Price: 108.49, UndiscountedPrice: 126.99, Stock: 27, Available: true, ImageUrl: "https://test.back.com/data/products/3/img1.png"},
			},
		},
		{
			name: "Ok without pagination",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "category_id", "name", "description", "price", "undiscounted_price", "stock", "available", "image_url"}).
					AddRow(1, 1, "product name 1", "product description 1", 0.99, 1.29, 12, true, "https://test.back.com/data/products/1/img1.png").
					AddRow(2, 1, "product name 2", "product description 2", 1.99, 1.99, 2, true, "https://test.back.com/data/products/2/img1.png").
					AddRow(3, 2, "product name 3", "product description 3", 108.49, 126.99, 27, true, "https://test.back.com/data/products/3/img1.png")
				mock.ExpectQuery("SELECT (.+) FROM products").
					WithArgs(0, 0).WillReturnRows(rows)
			},
			input: args{0, 0},
			want: []domain.Product{
				{Id: 1, CategoryId: 1, Name: "product name 1", Description: "product description 1", Price: 0.99, UndiscountedPrice: 1.29, Stock: 12, Available: true, ImageUrl: "https://test.back.com/data/products/1/img1.png"},
				{Id: 2, CategoryId: 1, Name: "product name 2", Description: "product description 2", Price: 1.99, UndiscountedPrice: 1.99, Stock: 2, Available: true, ImageUrl: "https://test.back.com/data/products/2/img1.png"},
				{Id: 3, CategoryId: 2, Name: "product name 3", Description: "product description 3", Price: 108.49, UndiscountedPrice: 126.99, Stock: 27, Available: true, ImageUrl: "https://test.back.com/data/products/3/img1.png"},
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

func TestGetProductById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProductPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		id int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    domain.Product
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "category_id", "name", "description", "price", "undiscounted_price", "stock", "available", "image_url"}).
					AddRow(1, 1, "product name 1", "product description 1", 0.99, 1.29, 12, true, "https://test.back.com/data/products/1/img1.png")
				mock.ExpectQuery("SELECT (.+) FROM products").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{1},
			want: domain.Product{
				Id:                1,
				CategoryId:        1,
				Name:              "product name 1",
				Description:       "product description 1",
				Price:             0.99,
				UndiscountedPrice: 1.29,
				Stock:             12,
				Available:         true,
				ImageUrl:          "https://test.back.com/data/products/1/img1.png",
			},
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "category_id", "name", "description", "price", "undiscounted_price", "stock", "available", "image_url"})

				mock.ExpectQuery("SELECT (.+) FROM products").
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

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProductPostgres(sqlx.NewDb(db, "sqlmock"), s)

	type args struct {
		input domain.CreateProductInput
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
				mock.ExpectQuery("INSERT INTO products").
					WithArgs(1, "Product name", "Product description", "", true, roundPrice(0.99), roundPrice(1.29), 12).WillReturnRows(rows)

				mock.ExpectExec("UPDATE products").
					WithArgs("https://test.back.com/data/products/1/img1.png", 1).WillReturnResult(driver.ResultNoRows)

				mock.ExpectCommit()
			},
			input: args{
				file: testFile,
				input: domain.CreateProductInput{
					CategoryId:        1,
					Name:              "Product name",
					Description:       "Product description",
					Price:             roundPrice(0.99),
					UndiscountedPrice: roundPrice(1.29),
					Stock:             12,
					Available:         true,
					ImgFile:           testFileHeader,
				},
			},
			want: 1,
		},
		{
			name: "CategoryNotFound",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO products").
					WithArgs(0, "Product name", "Product description", "", true, roundPrice(0.99), roundPrice(1.29), 12).WillReturnError(errors.New("category already exists"))
				mock.ExpectRollback()
			},
			input: args{
				file: testFile,
				input: domain.CreateProductInput{
					CategoryId:        0,
					Name:              "Product name",
					Description:       "Product description",
					Price:             roundPrice(0.99),
					UndiscountedPrice: roundPrice(1.29),
					Stock:             12,
					Available:         true,
					ImgFile:           testFileHeader,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateProduct(tt.input.input, tt.input.file)
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

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	fsTest := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: "https://test.back.com",
	})
	s := storage.NewStorage(fsTest)

	r := newProductPostgres(sqlx.NewDb(db, "sqlmock"), s)
	type args struct {
		productId int
		input     domain.UpdateProductInput
		file      multipart.File
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
				mock.ExpectQuery("UPDATE products SET (.+)").
					WithArgs(1, "new name", "new description", 12, roundPrice(0.99), roundPrice(1.29), true, "https://test.back.com/data/products/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				productId: 1,
				input: domain.UpdateProductInput{
					CategoryId:        intPointer(1),
					Name:              stringPointer("new name"),
					Description:       stringPointer("new description"),
					Stock:             intPointer(12),
					Price:             float32Pointer(roundPrice(0.99)),
					UndiscountedPrice: float32Pointer(roundPrice(1.29)),
					Available:         boolPointer(true),
					ImgFile:           testFileHeader,
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
				mock.ExpectQuery("UPDATE products SET (.+)").
					WithArgs(1, "new name", "new description", 12, roundPrice(0.99), roundPrice(1.29), "https://test.back.com/data/products/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				productId: 1,
				input: domain.UpdateProductInput{
					CategoryId:        intPointer(1),
					Name:              stringPointer("new name"),
					Description:       stringPointer("new description"),
					Stock:             intPointer(12),
					Price:             float32Pointer(roundPrice(0.99)),
					UndiscountedPrice: float32Pointer(roundPrice(1.29)),
					ImgFile:           testFileHeader,
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
				mock.ExpectQuery("UPDATE products SET (.+)").
					WithArgs(1, "new name", 12, roundPrice(0.99), roundPrice(1.29), "https://test.back.com/data/products/1/img1.png", 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				productId: 1,
				input: domain.UpdateProductInput{
					CategoryId:        intPointer(1),
					Name:              stringPointer("new name"),
					Stock:             intPointer(12),
					Price:             float32Pointer(roundPrice(0.99)),
					UndiscountedPrice: float32Pointer(roundPrice(1.29)),
					ImgFile:           testFileHeader,
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
				mock.ExpectQuery("UPDATE products SET (.+)").
					WithArgs(1, "new name", 12, roundPrice(0.99), roundPrice(1.29), true, 1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			input: args{
				productId: 1,
				input: domain.UpdateProductInput{
					CategoryId:        intPointer(1),
					Name:              stringPointer("new name"),
					Stock:             intPointer(12),
					Price:             float32Pointer(roundPrice(0.99)),
					UndiscountedPrice: float32Pointer(roundPrice(1.29)),
					Available:         boolPointer(true),
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateProduct(tt.input.productId, tt.input.input, tt.input.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func roundPrice(price float32) float32 {
	return float32(math.Round(float64(price)*100) / 100)
}

func intPointer(i int) *int {
	return &i
}

func float32Pointer(f float32) *float32 {
	return &f
}
