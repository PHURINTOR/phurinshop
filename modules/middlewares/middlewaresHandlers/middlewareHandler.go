package middlewaresHandlers

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewareUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// -------------------------------------------- Interface --------------------------
type IMidlewareHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
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
