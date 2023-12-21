package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mercan/ecommerce/internal/repositories/mongodb"
	"github.com/mercan/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongoRepository = mongodb.NewUserMongoRepository()

func IsEmailVerified(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(primitive.ObjectID)

	isVerified, err := mongoRepository.CheckEmailVerified(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.FailureAuthResponse{
			Message: "Internal Server Error",
			Status:  "error",
			Code:    fiber.StatusInternalServerError,
		})
	}

	if !isVerified {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.FailureAuthResponse{
			Message: "You can't access this resource. Please verify your email address.",
			Status:  "error",
			Code:    fiber.StatusUnauthorized,
		})
	}

	return ctx.Next()
}
