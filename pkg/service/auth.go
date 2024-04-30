package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
	"github.com/spf13/viper"
)

const confirmEmailTokenTTL = time.Hour
const signInTokenTTL = 24 * time.Hour

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
	confirmationToken, err := generateToken(tokenPayload, confirmEmailTokenTTL)
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
	tokenPayload, err := parseToken(token)

	if err != nil {
		return 0, err
	}
	user := new(domain.User)
	user.Name = tokenPayload["name"].(string)
	user.Email = tokenPayload["email"].(string)
	user.Password = generatePasswordHash(password)
	fmt.Println(user)
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
	return generateToken(payload, signInTokenTTL)

}

func (s *AuthService) ParseAuthToken(authToken string) (int, error) {
	tokenPayload, err := parseToken(authToken)
	if err != nil {
		return 0, err
	}

	id := tokenPayload["id"].(int)

	return id, nil
}

func generateToken(payload map[string]interface{}, tokenTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signKey := os.Getenv("TOKENS_SIGN_KEY")
	return token.SignedString([]byte(signKey))
}

func parseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKENS_SIGN_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload := claims["payload"].(map[string]interface{})
		return payload, nil
	}

	return nil, errors.New("invalid token")
}

func generatePasswordHash(password string) string {
	salt := os.Getenv("PASSWORD_HASH_SALT")
	hash := sha1.New()
	hash.Write([]byte(password))
	fmt.Print("SALTASAS")
	fmt.Print(salt)
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
