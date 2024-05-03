package service

import (
	"fmt"
	"os"

	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
	"github.com/spf13/viper"
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

func (s *ProfileService) RecoveryPassword(userId int) error {
	tokenPayload := map[string]interface{}{
		"id": userId,
	}
	signKey := os.Getenv("TOKEN_PASSWORD_RECOVERY_KEY")

	confirmationToken, err := generateToken(tokenPayload, signKey, confirmationTokenTTL)
	if err != nil {
		return err
	}

	user, err := s.repo.GetProfile(userId)
	if err != nil {
		return err
	}

	clientUrl := viper.GetString("client.password_recovery_url")
	confirmationLink := fmt.Sprintf("%s?confToken=%s", clientUrl, confirmationToken)

	const emailSubject = "Password recovery confirmation"
	emailBody := fmt.Sprintf("Please, enter through this link to change your password: %s", confirmationLink)

	return sendMail([]string{user.Email}, emailSubject, emailBody)
}

func (s *ProfileService) UpdatePassword(token, password string) error {
	signKey := os.Getenv("TOKEN_PASSWORD_RECOVERY_KEY")

	tokenPayload, err := parseToken(token, signKey)

	if err != nil {
		return err
	}

	userId := int(tokenPayload["id"].(float64))

	return s.repo.UpdatePassword(userId, generatePasswordHash(password))
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
