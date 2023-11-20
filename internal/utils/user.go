package utils

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/mercan/ecommerce/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

const (
	letterBytes               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	forgotPasswordTokenLength = 64
)

// GenerateVerificationCode generates a random 6 digit number
func GenerateVerificationCode() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 100000 and 999999 (6 digits) and convert it to string
	return strconv.Itoa(random.Intn(999999) + 100000)
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
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

// GenerateJWT generates a JWT token
func GenerateJWT(id primitive.ObjectID, role string) (string, error) {
	if role != "user" && role != "store" {
		return "", errors.New("invalid role")
	}

	token := jwt.New(jwt.SigningMethodHS512)
	now := time.Now().UTC()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(config.GetJWTConfig().Expiration * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["id"] = id
	claims["role"] = role
	claims["authorized"] = true

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(config.GetJWTConfig().Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ContextWithTimeout returns a context with a timeout
func ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout*time.Second)
}
