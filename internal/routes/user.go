package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mercan/ecommerce/internal/controllers"
	"github.com/mercan/ecommerce/internal/middleware"
	"github.com/mercan/ecommerce/internal/types"
	"time"
)

// SetupUserRoutes sets up user routes
func SetupUserRoutes(app *fiber.App) {
	userController := controllers.NewUserController()

	// Auth Group
	user := app.Group("/auth")

	limiterMiddleware := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(types.BaseResponse{
				Success: false,
				Error:   "Too many requests",
			})
		},
	})

	user.Use(limiterMiddleware)

	user.Post("/register", middleware.CheckContentType, userController.Register)
	user.Post("/login", middleware.CheckContentType, userController.Login)
	user.Get("/logout", middleware.IsAuthenticated, userController.Logout)

	user.Post("/change-password", middleware.CheckContentType, middleware.IsAuthenticated, middleware.IsEmailVerified, userController.ChangePassword)
	user.Post("/change-email", middleware.CheckContentType, middleware.IsAuthenticated, userController.ChangeEmail)

	user.Get("/verify-email", middleware.IsAuthenticated, userController.VerifyEmail)
	user.Get("/resend-verification-email", middleware.IsAuthenticated, userController.ResendEmailVerification)
	user.Get("/verify-phone", middleware.IsAuthenticated, userController.VerifyPhone)
	user.Get("/resend-verification-phone", middleware.IsAuthenticated, userController.ResendPhoneVerification)

	user.Get("/forgot-password", userController.ForgotPassword)
	//auth.Post("/refresh", userAuthController.Refresh)

}
