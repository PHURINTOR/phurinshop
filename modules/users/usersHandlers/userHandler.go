package usersHandlers

import (
	"strings"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/users"
	"github.com/PHURINTOR/phurinshop/modules/users/usersUsecases"
	"github.com/PHURINTOR/phurinshop/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

// -------------------------------------- Interface ---------------------------------
type IUserHandler interface {
	//User normal
	SignUpCustomer(c *fiber.Ctx) error

	//Login
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error

	//Logout
	SignOut(c *fiber.Ctx) error

	//Admin
	SignUpAdmin(c *fiber.Ctx) error        //SignUpAdmin
	GenerateAdminToken(c *fiber.Ctx) error //GenAdminkey

	//Check and Get
	GetUserProfile(c *fiber.Ctx) error
}

// -------------------------------------- Struct     ---------------------------------
// Enum error code
type usershandlerErrorCode string

const (
	signUpCustomer    usershandlerErrorCode = "users-001" // Error code SignUP Customer
	signInErr         usershandlerErrorCode = "users-002" // Error code signin
	refreshPassport1  usershandlerErrorCode = "users-003" // Error code signin
	signOutError      usershandlerErrorCode = "users-004" // Error code signin
	signUpAdmin       usershandlerErrorCode = "users-005" // Error code SignUP Admin
	generateAdminErr  usershandlerErrorCode = "users-006" // Error code SignUP Admin
	getUserProfileErr usershandlerErrorCode = "users-007" // Error code GetUserProfile

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

// ========================================= Method missing Func  =========================================
//=================================== Insert

// ---------------------- User
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

// ------------------- Admin
// *************************************  สิ่งที่ admin แตกต่างคือ  เวลาจะ Create Admin ต้องผ่าน middleware เพื่อใช้ key admin ประกอบเพื่อความปลอดภัย
func (h *userHandler) SignUpAdmin(c *fiber.Ctx) error {
	// Request body Parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdmin),
			err.Error()).Res()
	}

	//Email Validation
	if !req.IsEmail() {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdmin),
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

//	------------>  นำ register ไปใช้ใน module
//
// Genkey Admin
func (h *userHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewphurinshopAuth(
		auth.Admin,
		h.cfg.Jwt(),
		nil, //ไม่จำเป็นต้องมี เคลมเพราะไม่อยากให้มี payload
	)

	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateAdminErr),
			err.Error(),
		).Res()
	}

	return entities.NewErrorResponse(c).Success(
		fiber.StatusOK,
		&struct { //สร้าง Struct และกำหนดค่าเพื่อ return ทันที
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

// -------------------------------------- Authentication -----------------
func (h *userHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	//ข้างบนคือ body ผ่านมาได้แล้ว
	passport, err := h.userUsecase.GetPassport(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, passport).Res()
}

// RefreshPassport
func (h *userHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassport1),
			err.Error(),
		).Res()
	}
	//ข้างบนคือ body ผ่านมาได้แล้ว
	passport, err := h.userUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassport1),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, passport).Res()
}

// ----------------- Delete Oauth
func (h *userHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutError),
			err.Error(),
		).Res()
	}
	if err := h.userUsecase.DeleteOAuth(req.OauthId); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutError),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, nil).Res()
}

// ----------------- GetUserProfile
func (h *userHandler) GetUserProfile(c *fiber.Ctx) error {
	//Set Params
	userId := strings.Trim(c.Params("user_id"), " ") //ตัดคำช่องว่าง และ user_id คือ :user_id ใน path module

	//Get Profile
	result, err := h.userUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewErrorResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewErrorResponse(c).Success(fiber.StatusOK, result).Res()
}
