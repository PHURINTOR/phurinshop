package odersUsecases

import (
	"fmt"
	"math"

	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/orders"
	"github.com/PHURINTOR/phurinshop/modules/orders/odersRepositories"
	"github.com/PHURINTOR/phurinshop/modules/products/productsRepositories"
)

// ======================================= Interface =========================================
type IOrdersUsecase interface {
	FindOneOrder(orderId string) (*orders.Oders, error)
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Oders) (*orders.Oders, error)
	UpdateOrder(req *orders.Oders) (*orders.Oders, error)
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
// ---------- FindOneOrder -------------
func (u *odersUsecase) FindOneOrder(orderId string) (*orders.Oders, error) {
	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// ---------- FindManyOrder -------------
func (u *odersUsecase) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.ordersRepository.FindOrder(req)
	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

// ---------- Insert Order -------------
func (u *odersUsecase) InsertOrder(req *orders.Oders) (*orders.Oders, error) {
	// Check if products is exiets
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}
		prod, err := u.productsRepository.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}

		// Set price
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.ordersRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// ---------- Update Order ------------
func (u *odersUsecase) UpdateOrder(req *orders.Oders) (*orders.Oders, error) {

	if err := u.ordersRepository.UpdateOrder(req); err != nil {
		return nil, err // result not work
	}

	// result == work
	order, err := u.ordersRepository.FindOneOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return order, nil
}
