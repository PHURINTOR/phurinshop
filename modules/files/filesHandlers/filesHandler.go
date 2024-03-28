package filesHandlers

import (
	"github.com/PHURINTOR/phurinshop/config"
	fileUsecases "github.com/PHURINTOR/phurinshop/modules/files/FilesUsecases"
)

// ======================================= Interface =========================================
type IFilesHandler interface {
}

// ======================================= Struct ============================================
type filesHandler struct {
	cfg          config.IConfig
	fileUsecases fileUsecases.IFilesUsecase
}

// ======================================= Constructor =======================================
func FilesHandler(cfg config.IConfig, fileUsecases fileUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		fileUsecases: fileUsecases,
	}
}
