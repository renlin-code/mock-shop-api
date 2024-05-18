package service

import (
	"mime/multipart"

	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type CategoryService struct {
	repo repository.Category
}

func newCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo}
}

func (s *CategoryService) GetAll(limit, offset int) ([]domain.Category, error) {
	categories, err := s.repo.GetAll(limit, offset)
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return categories, errors_handler.NotFound("categories")
	}
	return categories, err
}

func (s *CategoryService) GetById(id int) (domain.Category, error) {
	category, err := s.repo.GetById(id)
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return category, errors_handler.NotFound("category")
	}
	return category, err
}

func (s *CategoryService) GetProducts(categoryId, limit, offset int) ([]domain.Product, error) {
	products, err := s.repo.GetProducts(categoryId, limit, offset)
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return products, errors_handler.NotFound("products")
	}
	return products, err
}

func (s *CategoryService) CreateCategory(input domain.CreateCategoryInput, file multipart.File) (int, error) {
	id, err := s.repo.CreateCategory(input, file)

	if errors_handler.ErrorIsType(err, errors_handler.TypeAlreadyExists) {
		return id, errors_handler.BadRequest("category with such name already exists")
	}
	return id, err
}

func (s *CategoryService) UpdateCategory(id int, input domain.UpdateCategoryInput, file multipart.File) error {
	err := s.repo.UpdateCategory(id, input, file)

	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return errors_handler.NotFound("category")
	}
	return err
}
