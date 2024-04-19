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
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		return []byte(config.GetJWTConfig().Secret), nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		// Check if token is in redis (blacklist)
		existingToken, err := authRedisRepo.IsTokenInBlacklist(token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(types.BaseResponse{
				Success: false,
				Error:   "Internal server error",
			})
		}

		if existingToken {
			return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
				Success: false,
				Error:   "Unauthorized",
			})
		}

		// Convert id to primitive.ObjectID
		userId, _ := primitive.ObjectIDFromHex(claims["id"].(string))
		// Set user id to context
		ctx.Locals("userId", userId)
		// Set user role to context
		ctx.Locals("role", claims["role"])
		// Set token expiration time to context
		ctx.Locals("exp", claims["exp"])
		// Set token to context
		ctx.Locals("token", token)

		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
		Success: false,
		Error:   "Unauthorized",
	})
}
