package productshandlers

import (
	"fmt"
	"strings"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/appinfo"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/files"
	"github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	"github.com/PHURINTOR/phurinshop/modules/products"
	productsusecases "github.com/PHURINTOR/phurinshop/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

// ======================================= Emums =======================================
type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "product-001"
	findProductsErr   productsHandlersErrCode = "product-002"
	InsertProductErr  productsHandlersErrCode = "product-003"
	UpdateProductErr  productsHandlersErrCode = "product-004"
	deleteProductErr  productsHandlersErrCode = "product-005"
)

// ======================================= Interface =========================================
type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProducts(c *fiber.Ctx) error
	AddProducts(c *fiber.Ctx) error
	UpdateProducts(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

// ======================================= Struct ============================================
type productsHandle struct {
	cfg             config.IConfig
	productsUsecase productsusecases.IProductsUsecase
	filesUsecases   filesUsecases.IFilesUsecase
}

// ======================================= Constructor =======================================
func ProductsHandle(cfg config.IConfig, productsUsecase productsusecases.IProductsUsecase, filesUsecases filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandle{
		cfg:             cfg,
		productsUsecase: productsUsecase,
		filesUsecases:   filesUsecases,
	}
}

// ======================================= Missing Function =======================================
// ------------  FindOneProduct ---------------
func (h *productsHandle) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")
	fmt.Print(productId)
	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	// Ok Res
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, product).Res()
}

// ------------  FindOneProduct ---------------
func (h *productsHandle) FindProducts(c *fiber.Ctx) error {

	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}
	// รับ query เป็น Struct = queryParser

	if err := c.QueryParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductsErr),
			err.Error(),
		).Res()
	}

	// Check And Defult Qeury
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	// Orderby
	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	// Sort
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	// ------- Use
	products := h.productsUsecase.FindProducts(req)

	// Ok Res
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, products).Res()
}

// ------------  Add Product ---------------
func (h *productsHandle) AddProducts(c *fiber.Ctx) error {

	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Images, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertProductErr),
			err.Error(),
		).Res()
	}

	// check
	if req.Category.Id <= 0 {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertProductErr),
			"category id is invalid",
		).Res()
	}

	// Use call Add Product
	products, err := h.productsUsecase.AddProducts(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.StatusInternalServerError,
			string(InsertProductErr),
			err.Error(),
		).Res()
	}
	// Ok Res
	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, products).Res()
}

// ------------  Update Product ---------------
func (h *productsHandle) UpdateProducts(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Images:   make([]*entities.Images, 0),
		Category: &appinfo.Category{},
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateProductErr),
			err.Error(),
		).Res()
	}

	req.Id = productId

	product, err := h.productsUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(UpdateProductErr),
			err.Error(),
		).Res()
	}
	// Ok res

	return entities.NewErrorResponse(c).Success(fiber.StatusOK, product).Res()
}

// ------------  Delete Product ---------------
func (h *productsHandle) DeleteProduct(c *fiber.Ctx) error {
	// 1. Delete local, 2. Delete From GCP
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	// delete

	// --- Stack sql command
	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("Images/test/%s", p.FileName),
		})
	}

	// Excute Delete GCP
	if err := h.filesUsecases.DeleteFileGCP(deleteFileReq); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}
	// Excute Delete Database local
	if err := h.productsUsecase.DeleteProduct(productId); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewErrorResponse(c).Success(fiber.StatusOK, nil).Res()
}
