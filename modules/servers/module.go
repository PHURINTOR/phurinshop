package servers

import (
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoHandlers"
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoRepositories"
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoUsecases"
	"github.com/PHURINTOR/phurinshop/modules/files/filesHandlers"
	"github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewareUsecases"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresHandlers"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresRepositories"
	monitorHanders "github.com/PHURINTOR/phurinshop/modules/monitors/monitorHandlers"
	"github.com/PHURINTOR/phurinshop/modules/orders/odersHandlers"
	"github.com/PHURINTOR/phurinshop/modules/orders/odersRepositories"
	"github.com/PHURINTOR/phurinshop/modules/orders/odersUsecases"
	productshandlers "github.com/PHURINTOR/phurinshop/modules/products/productsHandlers"
	productsrepositories "github.com/PHURINTOR/phurinshop/modules/products/productsRepositories"
	"github.com/PHURINTOR/phurinshop/modules/products/productsUsecases"
	"github.com/PHURINTOR/phurinshop/modules/users/usersHandlers"
	"github.com/PHURINTOR/phurinshop/modules/users/usersRepositories"
	"github.com/PHURINTOR/phurinshop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

//Menu แบบ factory

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OdersModule()
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
	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)

	//Login
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)

	//Create Admin
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)
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

	//Gen API key
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(1, 2), handler.GenerateApiKey)

	//Find Category
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)

	//Add Category
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	//Delete Category
	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)

}

// ============================================================ FilesModule ===========================================
// ============================  /v1/Files/ =================================
func (m *moduleFactory) FilesModule() {

	usecase := filesUsecases.FilesUsecase(m.server.cfg)
	handler := filesHandlers.FilesHandler(m.server.cfg, usecase)
	router := m.router.Group("/files")

	// Upload Files
	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)

	// Delete Files
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)

	// *** เหตุผลที่ใช้ Patch เพราะสามารถเพิ่ม Body เข้าไปได้

}

// ============================================================ ProductsModule ===========================================
func (m *moduleFactory) ProductsModule() {
	filesUsecases := filesUsecases.FilesUsecase(m.server.cfg)

	repository := productsrepositories.ProductsRepository(m.server.db, m.server.cfg, filesUsecases)
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productshandlers.ProductsHandle(m.server.cfg, usecase, filesUsecases)

	router := m.router.Group("/products")

	// FindOneProduct
	router.Get("/:product_id", m.mid.ApiKeyAuth(), handler.FindOneProduct)
	// FindProducts
	router.Get("/", m.mid.ApiKeyAuth(), handler.FindProducts)
	// AddProduct
	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddProducts)
	// UpdateProduct
	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateProducts)
	// DELETE
	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteProduct)
}

// ============================================================ OdersModule ===========================================
func (m *moduleFactory) OdersModule() {
	filesUsecases := filesUsecases.FilesUsecase(m.server.cfg)
	productRepository := productsrepositories.ProductsRepository(m.server.db, m.server.cfg, filesUsecases)

	repository := odersRepositories.OdersRepository(m.server.db)
	usecase := odersUsecases.OdersUsecase(repository, productRepository)
	handler := odersHandlers.OdersHandler(m.server.cfg, usecase)

	router := m.router.Group("/orders")

	// FindOneProduct
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.FindOneOrder)

	// FindOrder (Admin)
	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.FindOrder)

	// Insert Order
	router.Post("/", m.mid.JwtAuth(), handler.InsertOrder)

	// Update Order
	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.UpdateOrder)
}
