package productshandlers

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	productsusecases "github.com/PHURINTOR/phurinshop/modules/products/productsUsecases"
)

// ======================================= Emums =======================================
type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "product-001"
)

// ======================================= Interface =========================================
type IProductsHandler interface {
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
