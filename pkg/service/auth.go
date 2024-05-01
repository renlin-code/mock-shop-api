package service

import (
	"fmt"
	"os"

	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
	"github.com/spf13/viper"
)

type AuthService struct {
	repo repository.Authorization
}

func newAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) UserSignUp(name, email string) error {
	tokenPayload := map[string]interface{}{
		"name":  name,
		"email": email,
	}
	signKey := os.Getenv("TOKEN_SIGNUP_KEY")

	confirmationToken, err := generateToken(tokenPayload, signKey, confirmationTokenTTL)
	if err != nil {
		return err
	}
	clientUrl := viper.GetString("client.confirmation_email_url")
	confirmationLink := fmt.Sprintf("%s?confToken=%s", clientUrl, confirmationToken)

	const emailSubject = "Sign up confirmation"
	emailBody := fmt.Sprintf("Please, enter through this link to confirm your email: %s", confirmationLink)

	return sendMail([]string{email}, emailSubject, emailBody)
}

func (s *AuthService) CreateUser(token, password string) (int, error) {
	signKey := os.Getenv("TOKEN_SIGNUP_KEY")

	tokenPayload, err := parseToken(token, signKey)

	if err != nil {
		return 0, err
	}
	user := new(domain.User)
	user.Name = tokenPayload["name"].(string)
	user.Email = tokenPayload["email"].(string)
	user.Password = generatePasswordHash(password)
	return s.repo.CreateUser(*user)
}

func (s *AuthService) GenerateAuthToken(email, password string) (string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))

	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"id": user.Id,
	}

	signKey := os.Getenv("TOKEN_SIGNIN_KEY")

	return generateToken(payload, signKey, signInTokenTTL)

}

func (s *AuthService) ParseAuthToken(authToken string) (int, error) {
	signKey := os.Getenv("TOKEN_SIGNIN_KEY")
	tokenPayload, err := parseToken(authToken, signKey)
	if err != nil {
		return 0, err
	}

	id := int(tokenPayload["id"].(float64))

	return id, nil
}
