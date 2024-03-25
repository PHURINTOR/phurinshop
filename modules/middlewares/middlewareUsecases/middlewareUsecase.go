package middlewareUsecases

import (
	"github.com/PHURINTOR/phurinshop/modules/middlewares"
	"github.com/PHURINTOR/phurinshop/modules/middlewares/middlewaresRepositories"
)

// -------------------------------------------- Interface --------------------------
type IMiidlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool //middleware Token
	FindRole() ([]*middlewares.Role, error)          //find Role for autu
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

// --------------- Middleware User Token
func (u *middlewareUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}

// -----------Find Role
func (u *middlewareUsecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.middlewareRepository.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
