package productsUsecases

import productsrepositories "github.com/PHURINTOR/phurinshop/modules/products/productsRepositories"

// ======================================= Interface =========================================
type IProductsUsecase interface {
}

// ======================================= Struct ============================================
type productsUsecase struct {
	productsRepository productsrepositories.IProductsRepository
}

// ======================================= Constructor =======================================
func ProductsUsecase(productsRepository productsrepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}
