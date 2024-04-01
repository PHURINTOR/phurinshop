package productsrepositories

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	"github.com/jmoiron/sqlx"
)

// ======================================= Interface =========================================
type IProductsRepository interface {
}

// ======================================= Struct ============================================
type productsRepository struct {
	db            *sqlx.DB
	cfg           config.IConfig
	filesUsecases filesUsecases.IFilesUsecase
}

// ======================================= Constructor =======================================
func ProductsRepository(db *sqlx.DB, cfg config.IConfig, fileUsecases filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:            db,
		cfg:           cfg,
		filesUsecases: fileUsecases,
	}
}
