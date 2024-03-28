package filesUsecases

import "github.com/PHURINTOR/phurinshop/config"

// ======================================= Interface =========================================
type IFilesUsecase interface {
}

// ======================================= Struct ============================================
type filesUsecase struct {
	cfg config.IConfig
}

// ======================================= Constructor =======================================
func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}
