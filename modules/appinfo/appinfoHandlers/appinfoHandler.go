package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/appinfo"
	"github.com/PHURINTOR/phurinshop/modules/appinfo/appinfoUsecases"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

// ======================================= Enum Error Code =========================================
type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlersErrCode = "appinfo-002"
	AddCategoryErr    appinfoHandlersErrCode = "appinfo-003"
	RemoveCategoryErr appinfoHandlersErrCode = "appinfo-004"
)

// ======================================= Interface =========================================
type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error // Gen APIKey
	FindCategory(c *fiber.Ctx) error   // Find Category
	AddCategory(c *fiber.Ctx) error    // Add Category
	RemoveCategory(c *fiber.Ctx) error // Remove Category
}

// ======================================= Struct ============================================
type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

// ======================================= Constructor =======================================
func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

// =======================================Missing Function =======================================

// -------------------------------------- GenAPI-Key --------------------------------------
// Gen API key เพื่อเอาไปใช้ตรวจสอบใน middleware
func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewphurinshopAuth(
		auth.ApiKey,
		h.cfg.Jwt(),
		nil, //ไม่จำเป็นต้องมี paylaod
	)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

// -------------------------------------- Find Category --------------------------------------
func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {

	/*Note: queryParams = http://localhost:3000/v1/appinfo/category?title="1234"&b="12354" */

	req := new(appinfo.CategoryFilter) // alloate ไว้แล้วไม่ต้องผ่านแบบ pointer

	/*Trick 1 : if err := c.Query("category")   = query ตัวเดียว*/

	/*Trick 2 : req := appinfo.CategorFilter{}  = if err := c.QueryParser( &req); ***** ต้อง & เพราะไม่ไม่ได้ประกาศแบบ alloate ไว้ */

	if err := c.QueryParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}
	category, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, category).Res()
}

// -------------------------------------- Add Category --------------------------------------
func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddCategoryErr),
			err.Error(),
		).Res()
	}

	//ถ้าที่ส่งมามี 0 จะไม่ให้สร้าาง
	if len(req) == 0 {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			"categories request are empty",
		).Res()
	}

	//Use Add Category
	if err := h.appinfoUsecase.InsertCategory(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(AddCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, req).Res()
}

// -------------------------------------- Delete Category --------------------------------------
func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {
	// Input Parameter from path Param
	categoryId := strings.Trim(c.Params("category_id"), " ")

	// Convert type value Param to Int
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RemoveCategoryErr),
			"id type is invalid",
		).Res()
	}

	// Check Value Param
	if categoryIdInt <= 0 {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RemoveCategoryErr),
			"id must more than 0",
		).Res()
	}

	// Use Delete Category Process
	if err := h.appinfoUsecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(RemoveCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
