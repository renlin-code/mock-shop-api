package service

import (
	"fmt"
	"os"

	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
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
	user, err := s.repo.GetUserByEmail(email)
	if user.Id != 0 {
		return errors_handler.Forbidden("user with this email already exists")
	}
	if err != nil && !errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return err
	}

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
		return 0, errors_handler.BadRequest("invalid token")
	}
	var user domain.User
	user.Name = tokenPayload["name"].(string)
	user.Email = tokenPayload["email"].(string)
	user.Password = generatePasswordHash(password)

	id, err := s.repo.CreateUser(user)
	if err != nil {
		if errors_handler.ErrorIsType(err, errors_handler.TypeAlreadyExists) {
			return 0, errors_handler.Forbidden(err.Error())
		}
		return 0, err
	}

	return id, nil
}

func (s *AuthService) GenerateAuthToken(email, password string) (string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))

	if err != nil {
		if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
			return "", errors_handler.BadRequest("incorrect email or password")
		}
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
		return 0, errors_handler.BadRequest("invalid token")
	}

	id := int(tokenPayload["id"].(float64))

	return id, nil
}

func (s *AuthService) RecoveryPassword(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
			return errors_handler.NotFound("user with this email")
		}
		return err
	}

	tokenPayload := map[string]interface{}{
		"id": user.Id,
	}
	signKey := os.Getenv("TOKEN_PASSWORD_RECOVERY_KEY")

	confirmationToken, err := generateToken(tokenPayload, signKey, confirmationTokenTTL)
	if err != nil {
		return err
	}

	clientUrl := viper.GetString("client.password_recovery_url")
	confirmationLink := fmt.Sprintf("%s?confToken=%s", clientUrl, confirmationToken)

	const emailSubject = "Password recovery confirmation"
	emailBody := fmt.Sprintf("Please, enter through this link to change your password: %s", confirmationLink)

	return sendMail([]string{user.Email}, emailSubject, emailBody)
}

func (s *AuthService) UpdatePassword(token, password string) error {
	signKey := os.Getenv("TOKEN_PASSWORD_RECOVERY_KEY")

	tokenPayload, err := parseToken(token, signKey)

	if err != nil {
		return errors_handler.BadRequest("invalid token")
	}

	userId := int(tokenPayload["id"].(float64))

	err = s.repo.UpdatePassword(userId, generatePasswordHash(password))
	if err != nil && errors_handler.ErrorIsType(err, errors_handler.TypeNoRows) {
		return errors_handler.NotFound("user with this email")
	}
	return err
}
