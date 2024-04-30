package service

import (
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type ProductService struct {
	repo repository.Product
}

func newProductService(repo repository.Product) *ProductService {
	return &ProductService{repo}
}

func (s *ProductService) GetAll() ([]domain.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetById(id int) (domain.Product, error) {
	return s.repo.GetById(id)
}
