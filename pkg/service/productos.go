package service

import (
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type ProductService struct {
	repo repository.Product
}

func newProductService(repo repository.Product) *ProductService {
	return &ProductService{repo}
}

func (s *ProductService) GetAll() ([]domain.Product, error) {
	products, err := s.repo.GetAll()
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return products, errors_handler.NotFound("products")
	}
	return products, err
}

func (s *ProductService) GetById(id int) (domain.Product, error) {
	product, err := s.repo.GetById(id)
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return product, errors_handler.NotFound("product")
	}
	return product, err
}
