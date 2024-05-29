package service

import (
	"mime/multipart"

	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
)

type ProfileService struct {
	repo repository.Profile
}

func newProfileService(repo repository.Profile) *ProfileService {
	return &ProfileService{repo}
}

func (s *ProfileService) GetProfile(userId int) (domain.User, error) {
	return s.repo.GetProfile(userId)
}

func (s *ProfileService) UpdateProfile(userId int, input domain.UpdateProfileInput, file multipart.File) error {
	return s.repo.UpdateProfile(userId, input, file)
}

func (s *ProfileService) CreateOrder(userId int, products []domain.CreateOrderInputProduct) (int, error) {
	id, err := s.repo.CreateOrder(userId, products)
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return id, errors_handler.NotFound("product")
	}
	if errors_handler.ErrorIsType(err, errors_handler.TypeConstrainViolation) {
		return id, errors_handler.BadRequest("quantity exceeds the stock")
	}
	return id, err

}

func (s *ProfileService) GetAllOrders(userId, limit, offset int) ([]domain.Order, error) {
	return s.repo.GetAllOrders(userId, limit, offset)
}

func (s *ProfileService) GetOrderById(userId, orderId int) (domain.Order, error) {
	order, err := s.repo.GetOrderById(userId, orderId)
	if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return order, errors_handler.NotFound("order")
	}
	return order, err
}
