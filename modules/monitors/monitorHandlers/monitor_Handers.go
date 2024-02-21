package monitorHanders

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/monitors"
	"github.com/gofiber/fiber/v2"
)

//รับ API request ต่างๆ

type IMonitorHander interface {
	HealthCheck(c *fiber.Ctx) error //missing func
}

// struct
type monitorHander struct {
	cfg config.IConfig //ไม่ต้องต่อ database เลยนำเข้าแค่ config
}

// inital
func MonitorHandler(cfg config.IConfig) IMonitorHander { //return เป็นแบบ object
	return &monitorHander{ //stuct ตัวนี้ต้องเข้าถึง interface บนได้
		cfg: cfg,
	}
}

// implement missing function healcheck
func (h *monitorHander) HealthCheck(c *fiber.Ctx) error {
	res := &monitors.Monitor{
		Name:    h.cfg.App().Name(),
		Version: h.cfg.App().Version(),
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

//Export Handler ไปทำงานกับ module
