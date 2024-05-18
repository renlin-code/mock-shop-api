package repository

import (
	"mime/multipart"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email, password string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	UpdatePassword(userId int, password string) error
}

type Category interface {
	GetAll(limit, offset int) ([]domain.Category, error)
	GetById(id int) (domain.Category, error)
	GetProducts(categoryId, limit, offset int) ([]domain.Product, error)
	CreateCategory(input domain.CreateCategoryInput, file multipart.File) (int, error)
	UpdateCategory(id int, input domain.UpdateCategoryInput, file multipart.File) error
}

type Product interface {
	GetAll(limit, offset int) ([]domain.Product, error)
	GetById(id int) (domain.Product, error)
	CreateProduct(input domain.CreateProductInput, file multipart.File) (int, error)
	UpdateProduct(id int, input domain.UpdateProductInput, file multipart.File) error
}

type Profile interface {
	GetProfile(userId int) (domain.User, error)
	UpdateProfile(userId int, input domain.UpdateProfileInput, file multipart.File) error
	CreateOrder(userId int, products []domain.CreateOrderInputProduct) (int, error)
	GetAllOrders(userId, limit, offset int) ([]domain.Order, error)
	GetOrderById(userId, orderId int) (domain.Order, error)
}

type Repository struct {
	Authorization
	Category
	Product
	Profile
}

func NewRepository(db *sqlx.DB, s *storage.Storage) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
		Category:      newCategoryPostgres(db, s),
		Product:       newProductPostgres(db, s),
		Profile:       newProfilePostgres(db, s),
	}
}
