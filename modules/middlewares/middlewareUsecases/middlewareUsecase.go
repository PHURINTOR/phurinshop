package middlewareUsecases

import (
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresRepositories"
)

// -------------------------------------------- Interface --------------------------
type IMiidlewareUsecase interface {
}

// ------------------------------------ Struct ------------------------------------
type middlewareUsecase struct {
	middlewareRepository middlewaresRepositories.IMiidlewareRepository
}

// ------------------------------- Constructor -------------------------------------
func MiddlewareUsecase(middlewareRepository middlewaresRepositories.IMiidlewareRepository) IMiidlewareUsecase {
	return &middlewareUsecase{
		middlewareRepository: middlewareRepository,
	}
}
