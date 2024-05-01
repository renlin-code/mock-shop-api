package service

import (
	"crypto/sha1"
	"fmt"
	"os"
)

func generatePasswordHash(password string) string {
	salt := os.Getenv("PASSWORD_HASH_SALT")
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
