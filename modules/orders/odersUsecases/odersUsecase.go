package odersUsecases

import (
	"github.com/PHURINTOR/phurinshop/modules/orders/odersRepositories"
	"github.com/PHURINTOR/phurinshop/modules/products/productsRepositories"
)

// ======================================= Interface =========================================
type IOrdersUsecase interface {
}

// ======================================= Struct ============================================
type odersUsecase struct {
	ordersRepository   odersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

// ======================================= Constructor =======================================
func OdersUsecase(ordersRepository odersRepositories.IOrdersRepository, productsRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &odersUsecase{
		ordersRepository:   ordersRepository,
		productsRepository: productsRepository,
	}
}

// ======================================= missing Func =======================================
