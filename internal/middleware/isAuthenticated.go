package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/mercan/ecommerce/internal/types"
	"strings"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/repositories/redis"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var authRedisRepo = redis.NewAuthenticationRedisRepository()

func extractToken(Bearer string) (string, error) {
	splitToken := strings.Split(Bearer, "Bearer ")
	if len(splitToken) != 2 {
		return "", fiber.ErrUnauthorized
	}

	if splitToken[1] != "" {
		return splitToken[1], nil
	}

	return "", fiber.ErrUnauthorized
}

// IsAuthenticated middleware checks if the request has an Authorization header
func IsAuthenticated(ctx *fiber.Ctx) error {
	Bearer := ctx.Get("Authorization")
	if Bearer == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	token, err := extractToken(Bearer)
	if err != nil {
		return err
	}

	if existingToken, err := authRedisRepo.IsTokenInBlacklist(token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.BaseResponse{
			Success: false,
			Error:   "Internal server error",
		})
	} else if existingToken {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	// Parse and validate JWT token
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		return []byte(config.GetJWTConfig().Secret), nil
	})

	if err != nil || !jwtToken.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	userId, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	// Set user context with extracted data
	ctx.Locals("userId", userId)
	ctx.Locals("role", claims["role"])
	ctx.Locals("exp", claims["exp"])
	ctx.Locals("token", token)

	return ctx.Next()
}
