package helpers

import (
	"github.com/golang-jwt/jwt"
	"github.com/mercan/ecommerce/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// GenerateJWT generates a JWT token
func GenerateJWT(id primitive.ObjectID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	now := time.Now().UTC()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(config.GetJWTConfig().Expiration * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["id"] = id
	claims["authorized"] = true

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(config.GetJWTConfig().Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
