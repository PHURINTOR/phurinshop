package odersHandlers

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/orders/odersUsecases"
)

// ======================================= Interface =========================================
type IOdersHandler interface {
}

// ======================================= Struct ============================================
type odersHandler struct {
	cfg         config.IConfig
	oderUsecase odersUsecases.IOrdersUsecase
}

// ======================================= Constructor =======================================
func OdersHandler(cfg config.IConfig, oderUsecase odersUsecases.IOrdersUsecase) IOdersHandler {
	return &odersHandler{
		cfg:         cfg,
		oderUsecase: oderUsecase,
	}
}

// ======================================= missing Func =======================================
