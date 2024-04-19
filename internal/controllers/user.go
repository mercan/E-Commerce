package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mercan/ecommerce/internal/helpers"
	"github.com/mercan/ecommerce/internal/models"
	"github.com/mercan/ecommerce/internal/repositories/rabbitmq"
	"github.com/mercan/ecommerce/internal/services"
	"github.com/mercan/ecommerce/internal/types"
)

type UserController struct {
	userService services.UserService
	EmailQueue  rabbitmq.EmailQueueManager
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
		EmailQueue:  rabbitmq.NewEmailQueueManager(),
	}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	user := models.NewUser()

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	if err := user.RegisterValidation(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	token, err := controller.userService.Register(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Publish email verification message to RabbitMQ
	controller.EmailQueue.PublishEmailVerification(user.FirstName, user.Email)

	return ctx.Status(fiber.StatusCreated).JSON(types.UserRegisterResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
		Token: token,
	})
}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	var user models.UserLoginRequest

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	token, err := controller.userService.Login(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserLoginResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
		Token: token,
	})
}

func (controller *UserController) Logout(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	expFloat64 := ctx.Locals("exp").(float64)

	err := controller.userService.Logout(token, expFloat64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserLogoutResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
	})
}

func (controller *UserController) ChangePassword(ctx *fiber.Ctx) error {
	var user models.UserChangePasswordRequest

	token := ctx.Locals("token").(string)
	expFloat64 := ctx.Locals("exp").(float64)
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	if user.NewPassword != user.NewPasswordConfirm {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   "New password and new password confirm do not match",
		})
	}

	newToken, err := controller.userService.ChangePassword(userId, user, token, expFloat64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserChangePasswordResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
		Token: newToken,
	})
}

func (controller *UserController) ChangeEmail(ctx *fiber.Ctx) error {
	var user models.UserChangeEmailRequest

	token := ctx.Locals("token").(string)
	expFloat64 := ctx.Locals("exp").(float64)
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	newToken, err := controller.userService.ChangeEmail(userId, user, token, expFloat64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserChangeEmailResponse{
		BaseResponse: types.BaseResponse{
			Success: false,
		},
		Token: newToken,
	})
}

func (controller *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	var user models.UserVerificationRequest
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.QueryParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	err := controller.userService.VerifyEmail(userId, user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserVerifyEmailResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
	})
}

func (controller *UserController) ResendEmailVerification(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := controller.userService.ResendEmailVerification(userId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserResendEmailVerificationResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
	})
}

func (controller *UserController) VerifyPhone(ctx *fiber.Ctx) error {
	var user models.UserVerificationRequest
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.QueryParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	err := controller.userService.VerifyPhone(userId, user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserVerifyPhoneResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
	})

}

func (controller *UserController) ResendPhoneVerification(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := controller.userService.ResendPhoneVerification(userId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserResendPhoneVerificationResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
		},
	})
}

func (controller *UserController) ForgotPassword(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"token": helpers.GenerateForgotPasswordToken(),
	})
}
