package service

import (
	"mime/multipart"

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

func (s *ProductService) GetAll(limit, offset int, search string) ([]domain.Product, error) {
	products, err := s.repo.GetAll(limit, offset, search)
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

func (s *ProductService) GetFilePath(productId int, fileName string) string {
	return s.repo.GetFilePath(productId, fileName)
}

func (s *ProductService) CreateProduct(input domain.CreateProductInput, file multipart.File) (int, error) {
	id, err := s.repo.CreateProduct(input, file)

	if errors_handler.ErrorIsType(err, errors_handler.TypeForeignKeyViolation) {
		return id, errors_handler.BadRequest("provided category_id does not correspond to any existing category")
	}
	return id, err
}

func (s *ProductService) UpdateProduct(id int, input domain.UpdateProductInput, file multipart.File) error {
	err := s.repo.UpdateProduct(id, input, file)

	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return errors_handler.NotFound("product")
	}
	return err
}
