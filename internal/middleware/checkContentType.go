package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mercan/ecommerce/internal/utils"
)

// CheckContentType middleware checks if the request content type is application/json
func CheckContentType(ctx *fiber.Ctx) error {
	// Check if content type is application/json or application/json; charset=utf-8
	if ctx.Get("Content-Type") != fiber.MIMEApplicationJSON &&
		ctx.Get("Content-Type") != fiber.MIMEApplicationJSONCharsetUTF8 {
		return ctx.Status(fiber.StatusUnsupportedMediaType).JSON(utils.FailureAuthResponse{
			Message: "Content-Type must be application/json or application/json; charset=utf-8",
			Status:  "error",
			Code:    fiber.StatusUnsupportedMediaType,
		})
	}

	return ctx.Next()
}
