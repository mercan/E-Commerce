package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mercan/ecommerce/internal/models"
	"github.com/mercan/ecommerce/internal/repositories/rabbitmq"
	"github.com/mercan/ecommerce/internal/services"
	"github.com/mercan/ecommerce/internal/types"
	"github.com/mercan/ecommerce/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userService  services.IUserService
	rabbitMQRepo *rabbitmq.Repository
}

func NewUserController() *UserController {
	return &UserController{
		userService:  services.NewUserService(),
		rabbitMQRepo: rabbitmq.NewRepository(),
	}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	user := models.NewUser()

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.UserRegisterResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	ip := "85.108.1.177" // IP address for testing purposes only - ctx.IP() is not working on localhost
	user.LoginHistory = append(user.LoginHistory, models.LoginHistory{
		IP: ip,
	})
	userAgent := ctx.GetReqHeaders()["User-Agent"]

	token, err := controller.userService.Register(user, userAgent)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserRegisterResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	// Publish email verification message to RabbitMQ
	controller.rabbitMQRepo.PublishEmailVerification(user.FirstName, user.Email)

	return ctx.Status(fiber.StatusCreated).JSON(types.UserRegisterResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusCreated,
		},
		Token: token,
	})
}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	var user models.UserLoginInput

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.UserLoginResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	ip := "78.176.227.107" // IP address for testing purposes only - ctx.IP() is not working on localhost
	userAgent := ctx.GetReqHeaders()["User-Agent"]

	token, err := controller.userService.Login(user, ip, userAgent)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserLoginResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserLoginResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
		Token: token,
	})
}

func (controller *UserController) Logout(ctx *fiber.Ctx) error {
	token := ctx.Locals("token").(string)
	expFloat64 := ctx.Locals("exp").(float64)

	err := controller.userService.Logout(token, expFloat64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserLogoutResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserLogoutResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
	})
}

func (controller *UserController) ChangePassword(ctx *fiber.Ctx) error {
	var user models.UserChangePasswordInput

	token := ctx.Locals("token").(string)
	expFloat64 := ctx.Locals("exp").(float64)
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.UserChangePasswordResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	if user.NewPassword != user.NewPasswordConfirm {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserChangePasswordResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: "New password and new password confirm do not match",
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	newToken, err := controller.userService.ChangePassword(userId, user, token, expFloat64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserChangePasswordResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserChangePasswordResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
		Token: newToken,
	})
}

func (controller *UserController) ChangeEmail(ctx *fiber.Ctx) error {
	var user models.UserChangeEmailInput

	token := ctx.Locals("token").(string)
	expFloat64 := ctx.Locals("exp").(float64)
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.UserChangeEmailResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	newToken, err := controller.userService.ChangeEmail(userId, user, token, expFloat64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserChangeEmailResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserChangeEmailResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
		Token: newToken,
	})
}

func (controller *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	var user models.UserVerificationInput
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.QueryParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.UserVerifyEmailResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	err := controller.userService.VerifyEmail(userId, user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserVerifyEmailResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserVerifyEmailResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
	})
}

func (controller *UserController) ResendEmailVerification(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := controller.userService.ResendEmailVerification(userId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserResendEmailVerificationResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserResendEmailVerificationResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
	})
}

func (controller *UserController) VerifyPhone(ctx *fiber.Ctx) error {
	var user models.UserVerificationInput
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := ctx.QueryParser(&user); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(types.UserVerifyPhoneResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	err := controller.userService.VerifyPhone(userId, user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserVerifyPhoneResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusBadRequest,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserVerifyPhoneResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
	})

}

func (controller *UserController) ResendPhoneVerification(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(primitive.ObjectID)

	if err := controller.userService.ResendPhoneVerification(userId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.UserResendPhoneVerificationResponse{
			UserBaseResponse: types.UserBaseResponse{
				ErrorResponse: types.ErrorResponse{
					Error: err.Error(),
				},
				Code: fiber.StatusUnprocessableEntity,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.UserResendPhoneVerificationResponse{
		UserBaseResponse: types.UserBaseResponse{
			Code: fiber.StatusOK,
		},
	})
}

func (controller *UserController) ForgotPassword(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"token": utils.GenerateForgotPasswordToken(),
	})
}
