package servers

import (
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoHandlers"
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoRepositories"
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoUsecases"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewareUsecases"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresHandlers"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresRepositories"
	monitorHanders "github.com/PHURINTOR/phurinshop/modules/monitors/monitorHandlers"
	"github.com/PHURINTOR/phurinshop/modules/users/usersHandlers"
	"github.com/PHURINTOR/phurinshop/modules/users/usersRepositories"
	"github.com/PHURINTOR/phurinshop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

//Menu แบบ factory

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
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

// ============================================================ UserModule ===========================================
func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UserRepository(m.server.db)
	usecase := usersUsecases.UserUsecase(m.server.cfg, repository)
	handler := usersHandlers.UserHandler(m.server.cfg, usecase)

	//============================  /v1/users/ =================================
	router := m.router.Group("/users")

	//Create User
	router.Post("/signup", handler.SignUpCustomer)

	//Login
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)

	//Create Admin
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)
	// Initail Admin เข้ามาใน database 1 คน
	// generate Admin Key
	// ทุกครั้งที่ Create Admin เพิ่มต้องเอา Admin key Token มาด้วย ผ่าน middleware
	//** ต่อให้ใครได้ Admin Token ไป  แต่ไม่มี key Admin จาก Env มาประกอบก็จะไม่สามารถ Gen ได้
	// Newadmin = AdminToken + key Admin (env)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
	//router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(1, 2), handler.GenerateAdminToken)

	//Authorization
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile) //Path param
}

// ============================================================ AppinfoModule ===========================================
// ============================  /v1/Appinfo/ =================================
func (m *moduleFactory) AppinfoModule() {

	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)
	router := m.router.Group("/appinfo")
	_ = router
	_ = handler

}
