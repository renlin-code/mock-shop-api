package service

import (
	"mime/multipart"

	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type Authorization interface {
	UserSignUp(name, email string) error
	CreateUser(token, password string) (int, error)
	GenerateAuthToken(email, password string) (string, error)
	ParseAuthToken(token string) (int, error)
	RecoveryPassword(email string) error
	UpdatePassword(token, password string) error
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
	UpdateProfile(userId int, input domain.UpdateProfileInput, file multipart.File) error
	CreateOrder(userId int, products []domain.CreateOrderInputProduct) (int, error)
	GetAllOrders(userId int) ([]domain.Order, error)
	GetOrderById(userId, orderId int) (domain.Order, error)
}

type Service struct {
	Authorization
	Category
	Product
	Profile
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
		Category:      newCategoryService(repos.Category),
		Product:       newProductService(repos.Product),
		Profile:       newProfileService(repos.Profile),
	}
}
