package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// =========================================   Struct  =========================================
type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

// -------------------------------- Register
type UserRegisterReq struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
	Username string `db:"username" json:"username" form:"username"`
}

// -------------------------------- SignIn
type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type UserPassport struct {
	User  *User      `json:"user"` //response กลับตอน login
	Token *UserToken `json:"token"`
}

// -------------------------------- Authentication
type UserCredential struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
	RoleId   int    `db:"role_id"`
}

// -------------------------------- JWT
type UserClaims struct {
	//ไม่ควรเป็นข้อมูลที่ secert จนเกินไป เพราะว่า playload/Claims จะเข้ารหัสแบบ base 64 ต่อให้ไม่มี key ก็สามารถถอดรหัสได้
	Id     string `db:"id" json:"id"`
	RoleId int    `db:"roleid" json:"roleid"`
}

// -------------------------------- User Refresh Token
type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"` //Check Parse Token Refresh  ==> auth
}

// -------------------------------- User Oauth
type Oauth struct {
	Id     string `db:"id" json:"id"`           //Role Id
	UserId string `db:"user_id" json:"user_id"` //User Id
}

// -------------------------------- User Remove Oauth
type UserRemoveCredential struct {
	OauthId string `json:"oauth_id" form:"oauth_id"`
}

// =========================================  Function missing of Struct =========================================

// -------------------------------- Register Function
// Bcrypt Password
func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hashed password failed: %v", err)
	}
	obj.Password = string(hashedPassword)
	return nil
}

// Check email
func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}
