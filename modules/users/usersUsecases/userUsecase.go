package usersUsecases

import (
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/users"
	usersrepositories "github.com/PHURINTOR/phurinshop/modules/users/usersRepositories"
)

// ========================================== Interface ==========================================

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	//InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
}

// ========================================== Struct     ==========================================
type userUsecase struct {
	cfg             config.IConfig
	usersRepository usersrepositories.IUserRepository
}

// ========================================== Constructor ========================================
func UserUsecase(cfg config.IConfig, usersRepository usersrepositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

// ========================================== Method missing Func ========================================

// -------------------------------- Register Function
func (u *userUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {

	//ก่อนส่งข้อมูลเข้าไป hash password ก่อน
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//Insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}
