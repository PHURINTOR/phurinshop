package servers

import (
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
}

// constructor , inital
func NewModule(r fiber.Router, s *server) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
	}
}

// implement missing
func (m *moduleFactory) MonitorModule() { //ไม่มี return เพราะอยากทำให้เป็น router เฉยๆ  แล้ว export ออกไปใช้ใน server
	handler := monitorHanders.MonitorHandler(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}
