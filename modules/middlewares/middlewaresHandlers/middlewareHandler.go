package middlewaresHandlers

import (
	"strings"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewareUsecases"
	"github.com/PHURINTOR/phurinshop/pkg/auth"
	"github.com/PHURINTOR/phurinshop/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// -------------------------------------------- Interface --------------------------
type IMidlewareHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler

	//User Token
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler

	//Authorizetion
	Authorize(expectRoleId ...int) fiber.Handler

	//======================= ApiKey
	//ApiKeyAuth
	ApiKeyAuth() fiber.Handler
}

// ------------------------------------ Struct ------------------------------------
type middlewareHandlers struct {
	cfg                 config.IConfig
	middlewaresUsecases middlewareUsecases.IMiidlewareUsecase
}

// Enum error
type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	ParamsCheck    middlewareHandlersErrCode = "middleware-003"

	//find Role
	authorizeErr middlewareHandlersErrCode = "middleware-004"

	//============ ApiKey
	//ApikeyAuth
	apiKeyAuthErr middlewareHandlersErrCode = "middleware-005"
)

// ------------------------------- Constructor -------------------------------------
func MiddlewareHandler(cfg config.IConfig, middlewaresUsecases middlewareUsecases.IMiidlewareUsecase) IMidlewareHandlers {
	return &middlewareHandlers{
		cfg:                 cfg,
		middlewaresUsecases: middlewaresUsecases,
	}
}

// ------------------------------- implement missing -------------------------------------
func (h *middlewareHandlers) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",                              //เข้าถึงได้ทุก IP
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH", //เข้าถึงได้ทุก Method
		AllowHeaders:     "",
		AllowCredentials: false, //เดี๋ยวจะใช้ Token
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewareHandlers) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewareHandlers) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path} \n",
		TimeFormat: "01/02/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

// ----------------------------------- Middleware User Token
func (h *middlewareHandlers) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ") // ตัดคำ String,  prefix
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecases.FindAccessToken(claims.Id, token) {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

// Check user_id (Access Token)  =  user_id (Profile)
func (h *middlewareHandlers) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 { //เช็คว่าเป็น admin ไหม
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(ParamsCheck),
				"never gonna give you up",
			).Res()
		}
		return c.Next() //ทำงานเสร็จจะส่งต่อให้ middleware ต่อไป
	}
}

// ----------------------------------- Middleware Find Role-Based
func (h *middlewareHandlers) Authorize(expectRoleId ...int) fiber.Handler {
	//expect คือ role ที่คาดหวังว่าสามารถใช้ได้

	return func(c *fiber.Ctx) error {
		//AuthJWT run id role ==> Locals
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}
		roles, err := h.middlewaresUsecases.FindRole()
		if err != nil {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}

		//compair
		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewErrorResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission",
		).Res()
	}
}

// ----------------------------------- Middleware ApiKey
func (h *middlewareHandlers) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")
		if _, err := auth.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyAuthErr),
				"apikey is invalid or required",
			).Res()

		}
		return c.Next()
	}
}
