package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mercan/ecommerce/internal/repositories/mongodb"
	"github.com/mercan/ecommerce/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongoRepository = mongodb.NewUserMongoRepository()

func IsEmailVerified(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(primitive.ObjectID)

	isVerified, err := mongoRepository.CheckEmailVerified(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.BaseResponse{
			Success: false,
			Error:   "Internal server error",
		})
	}

	if !isVerified {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.BaseResponse{
			Success: false,
			Error:   "You can't access this resource. Please verify your email address.",
		})
	}

	return ctx.Next()
}
