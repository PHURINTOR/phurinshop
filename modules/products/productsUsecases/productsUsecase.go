package productsUsecases

import (
	"math"

	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/products"
	productsrepositories "github.com/PHURINTOR/phurinshop/modules/products/productsRepositories"
)

// ======================================= Interface =========================================
type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) *entities.PaginateRes
	AddProducts(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
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

// ======================================= Missing Function =======================================
// ------------  FindOneProduct ---------------
func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// ------------  FindProducts ---------------
func (u *productsUsecase) FindProducts(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productsRepository.FindProducts(req)

	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))), //หาญปัดเศษ = จำนวนทั้งหมด(count) / จำนวน limit
	}
}

// ------------  InsertProduct ---------------
func (u *productsUsecase) AddProducts(req *products.Product) (*products.Product, error) {
	products, err := u.productsRepository.InsertProducts(req)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ------------  UpdateProduct ---------------
func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	products, err := u.productsRepository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ------------  DeleteProdcut ---------------
func (u *productsUsecase) DeleteProduct(productId string) error {
	if err := u.productsRepository.DeleteProduct(productId); err != nil {
		return err
	}
	return nil
}
