package usersRepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/PHURINTOR/phurinshop/modules/users"
	"github.com/PHURINTOR/phurinshop/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

// -------------------------------------- Interface ---------------------------------
type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)

	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) // read hash for repair

	InsertOauth(req *users.UserPassport) error

	FindOneOauth(refreshToken string) (*users.Oauth, error) //refresh Token Gen
	UpdateOauth(req *users.UserToken) error                 //refresh Token Gen
	DeleteOauth(oauthId string) error

	//GetProfile
	GetProfile(userId string) (*users.User, error)
}

// -------------------------------------- Struct     ---------------------------------
type userRepository struct {
	db *sqlx.DB
}

// -------------------------------------- Constructor--------------------------------
func UserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

// ============================ missing Functin =================================

// ----------------------------------------------- Insert User
func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {

	result := usersPatterns.InsertUser(r.db, req, isAdmin)
	var err error

	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	//Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ----------------------------------------------- Authentication
// read hash for repair
func (r *userRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"username",
		"role_id"
	FROM "users"
	WHERE "email" = $1;
	`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		//สามารถ past query เข้าไปได้เลยเพราะว่า UserCredentialCheck struct สร้างไว้แล้ว ชื่อต้องตรง หรือใช้ as
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// ----------------------------------------------- Insert Oauth
func (r *userRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //เป็นการกำหนดระยะเวลาในการ insert
	defer cancel()

	//========================================== เพิ่ม Token ใน  table auth แล้ว Return
	query := `
	INSERT INTO "oauth" (
		"user_id",
		"access_token",
		"refresh_token"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.AccessToken,
		req.Token.RefreshToken,
	).Scan(&req.Token.Id); err != nil { //เราสามารถ scan เข้าไปได้เลย แต่ต้องเอา & มารองรับ
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

// ----------------------------------------------- Update Oauth
// *** จำเป็นต้องสร้าง RepeatToken ในไฟล์ Auth เสียก่อน
//  1. FindOneOauth  2. UpdateOauth

// FindOneOauth
// *** จำเป็นต้องสร้าง struct = user.Oauth ในไฟล์ users เสียก่อน
func (r *userRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	query := `
	SELECT
		"id",
		"user_id"
	FROM "oauth"
	WHERE "refresh_token" = $1;`
	oauth := new(users.Oauth)
	fmt.Println(oauth.Id, oauth.UserId)
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauth, nil
}

// UpdateOauth
func (r *userRepository) UpdateOauth(req *users.UserToken) error {
	query := `
	UPDATE "oauth" SET
		"access_token" = :access_token,
		"refresh_token" = :refresh_token
	WHERE "id"= :id;`

	/*:access_token คือ match ชื่อใน stuct db:"access_token" อัตโนมัติ   วิธีนี้นิยมใช้กับ Update เพราะไม่เจาะจง*/

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update Oauth Failed: %v", err)
	}
	// NamedExecContext หากไม่สำเร็จ จะ return resource อย่างปลอดภัย ไม่กระทบ
	return nil
}

// ----------------------  Delete OAuth
func (r *userRepository) DeleteOauth(oauthId string) error {
	query := `DELETE FROM "oauth" WHERE "id" = $1;`
	if _, err := r.db.ExecContext(context.Background(), query, oauthId); err != nil { //context คือ ถ้าคำสั่งนี้ไม่สามารถทำจบในเวลา
		return fmt.Errorf("oauth not found: %v", err)
	}
	return nil
}

// ----------------------  GetProfile
func (r *userRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"id",
		"email",
		"username",
		"role_id"
	FROM "users"
	WHERE "id"= $1;`
	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}
