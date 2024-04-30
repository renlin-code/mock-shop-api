package service

import (
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type CategoryService struct {
	repo repository.Category
}

func newCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo}
}

func (s *CategoryService) GetAll() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetById(id int) (domain.Category, error) {
	return s.repo.GetById(id)
}
