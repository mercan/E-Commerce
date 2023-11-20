package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"github.com/mercan/ecommerce/internal/config"
	"github.com/mercan/ecommerce/internal/repositories/redis"
	"github.com/mercan/ecommerce/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var redisRepository = redis.NewRepository()

// IsAuthenticated middleware checks if the request has an Authorization header
func IsAuthenticated(ctx *fiber.Ctx) error {
	Bearer := ctx.Get("Authorization")
	if Bearer == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.FailureAuthResponse{
			Message: "Unauthorized",
			Status:  "error",
			Code:    fiber.StatusUnauthorized,
		})
	}

	token := Bearer[7:]
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.FailureAuthResponse{
			Message: "Unauthorized",
			Status:  "error",
			Code:    fiber.StatusUnauthorized,
		})
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		return []byte(config.GetJWTConfig().Secret), nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.FailureAuthResponse{
			Message: "Unauthorized",
			Status:  "error",
			Code:    fiber.StatusUnauthorized,
		})
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		// Check if token is in redis (blacklist)
		existingToken, err := redisRepository.IsTokenInBlacklist(token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.FailureAuthResponse{
				Message: "Internal Server Error",
				Status:  "error",
				Code:    fiber.StatusInternalServerError,
			})
		}

		if existingToken {
			return ctx.Status(fiber.StatusUnauthorized).JSON(utils.FailureAuthResponse{
				Message: "Unauthorized",
				Status:  "error",
				Code:    fiber.StatusUnauthorized,
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

	return ctx.Status(fiber.StatusUnauthorized).JSON(utils.FailureAuthResponse{
		Message: "Unauthorized",
		Status:  "error",
		Code:    fiber.StatusUnauthorized,
	})
}
