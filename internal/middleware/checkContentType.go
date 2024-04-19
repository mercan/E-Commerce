package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mercan/ecommerce/internal/types"
)

// CheckContentType middleware checks if the request content type is application/json
func CheckContentType(ctx *fiber.Ctx) error {
	// Check if content type is application/json or application/json; charset=utf-8
	if ctx.Get("Content-Type") != fiber.MIMEApplicationJSON &&
		ctx.Get("Content-Type") != fiber.MIMEApplicationJSONCharsetUTF8 {
		return ctx.Status(fiber.StatusUnsupportedMediaType).JSON(types.BaseResponse{
			Success: false,
			Error:   "Content type must be application/json or application/json; charset=utf-8",
		})
	}

	return ctx.Next()
}
