package service

import (
	"github.com/renlin-code/mock-shop-api/pkg/domain"
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

func (s *ProfileService) UpdateProfile(userId int, input domain.UpdateProfileInput) error {
	return s.repo.UpdateProfile(userId, input)
}

func (s *ProfileService) CreateOrder(userId int, products []domain.CreateOrderInputProduct) (int, error) {
	return s.repo.CreateOrder(userId, products)
}

func (s *ProfileService) GetAllOrders(userId int) ([]domain.Order, error) {
	return s.repo.GetAllOrders(userId)
}

func (s *ProfileService) GetOrderById(userId, orderId int) (domain.Order, error) {
	return s.repo.GetOrderById(userId, orderId)
}
