package usersUsecases

import (
	"fmt"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/users"
	usersrepositories "github.com/PHURINTOR/phurinshop/modules/users/usersRepositories"
	"github.com/PHURINTOR/phurinshop/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

// ========================================== Interface ==========================================

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) //insert user
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)    //insert user

	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOAuth(oauthId string) error

	GetUserProfile(userid string) (*users.User, error) //Check Params Token

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

// Admin
func (u *userUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {

	//ก่อนส่งข้อมูลเข้าไป hash password ก่อน
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//Insert user
	result, err := u.usersRepository.InsertUser(req, true) //true เพราะเป็น Admin
	if err != nil {
		return nil, err
	}
	return result, nil
}

// -------------------------------- Authentication
func (u *userUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {

	//---------------------    Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//---------------------  Compair Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	accessToken, err := auth.NewphurinshopAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, err := auth.NewphurinshopAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	//Set passport //return token ตาม struct Userpassport

	passport := &users.UserPassport{ //&user คือ ประกาส sturct พร้อมใส่ค่า init ด้วย
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil { //insert Oauth
		return nil, err
	}
	return passport, nil
}

//--------------------- RefreshPassport

func (u *userUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {

	// Parse token  ดึง payload
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check oauth  ดึง payloadไปหา Oauth ว่าเคย login หรือไม่
	fmt.Printf("RefreshToken On Usecase := %v", req.RefreshToken)
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find profile เอา Oauth มาเช็คว่ามี profile จริงหรือไม่
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	//ทำ payload ใหม่
	newCliams := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewphurinshopAuth(
		auth.Access,
		u.cfg.Jwt(),
		newCliams,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(
		u.cfg.Jwt(),
		newCliams,
		claims.ExpiresAt.Unix(),
	)

	//make Passport payload จากที่สร้างใหม่
	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil

}

// ---------------------- delete Refresh Token
func (u *userUsecase) DeleteOAuth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

// -------------------------Get User Profile
func (u *userUsecase) GetUserProfile(userid string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userid)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
