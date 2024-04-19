package helpers

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	letterBytes               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	forgotPasswordTokenLength = 64
)

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// VerifyPassword compares a hashed password with a plain text password
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateForgotPasswordToken generates a random string
func GenerateForgotPasswordToken() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, forgotPasswordTokenLength)

	for i := range bytes {
		bytes[i] = letterBytes[random.Intn(len(letterBytes))]
	}

	return string(bytes)
}
