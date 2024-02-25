package servers

import (
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewareUsecases"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresHandlers"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresRepositories"
	monitorHanders "github.com/PHURINTOR/phurinshop/modules/monitors/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

//Menu แบบ factory

type IModuleFactory interface {
	MonitorModule()
}

//struct

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.IMidlewareHandlers
}

// constructor , inital
func NewModule(r fiber.Router, s *server, mid middlewaresHandlers.IMidlewareHandlers) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMidlewareHandlers {
	repository := middlewaresRepositories.MiddlewareRepository(s.db)
	usecase := middlewareUsecases.MiddlewareUsecase(repository)
	return middlewaresHandlers.MiddlewareHandler(s.cfg, usecase)
}

// implement missing
func (m *moduleFactory) MonitorModule() { //ไม่มี return เพราะอยากทำให้เป็น router เฉยๆ  แล้ว export ออกไปใช้ใน server
	handler := monitorHanders.MonitorHandler(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}
