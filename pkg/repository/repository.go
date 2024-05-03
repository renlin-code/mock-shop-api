package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email, password string) (domain.User, error)
}

type Category interface {
	GetAll() ([]domain.Category, error)
	GetById(id int) (domain.Category, error)
	GetProducts(categoryId int) ([]domain.Product, error)
}

type Product interface {
	GetAll() ([]domain.Product, error)
	GetById(id int) (domain.Product, error)
}

type Profile interface {
	GetProfile(userId int) (domain.User, error)
	UpdateProfile(userId int, input domain.UpdateProfileInput) error
	UpdatePassword(userId int, password string) error
	CreateOrder(userId int, products []domain.CreateOrderInputProduct) (int, error)
	GetAllOrders(userId int) ([]domain.Order, error)
	GetOrderById(userId, orderId int) (domain.Order, error)
}

type Repository struct {
	Authorization
	Category
	Product
	Profile
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
		Category:      newCategoryPostgres(db),
		Product:       newProductPostgres(db),
		Profile:       newProfilePostgres(db),
	}
}
