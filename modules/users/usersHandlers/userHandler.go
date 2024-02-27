package usersHandlers

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/users"
	"github.com/PHURINTOR/phurinshop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

// -------------------------------------- Interface ---------------------------------
type IUserHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

// -------------------------------------- Struct     ---------------------------------
// Enum error code
type usershandlerErrorCode string

const (
	signUpCustomer usershandlerErrorCode = "users-001"
)

type userHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

// -------------------------------------- Constructor--------------------------------
func UserHandler(cfg config.IConfig, userUsecase usersUsecases.IUserUsecase) IUserHandler {
	return &userHandler{
		cfg:         cfg,
		userUsecase: userUsecase,
	}
}

// --------------------------------------- Method missing Func --------------------------------
func (h *userHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body Parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomer),
			err.Error()).Res()
	}

	//Email Validation
	if !req.IsEmail() {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomer),
			"email pattern is invalid",
		).Res()
	}

	//Insert
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomer),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomer),
				err.Error(),
			).Res()
		default:
			return entities.NewErrorResponse(c).Error(
				fiber.ErrInternalServerError.Code, //500
				string(signUpCustomer),
				err.Error(),
			).Res()
		}
	}

	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, result).Res()
}

//นำ register ไปใช้ใน module
