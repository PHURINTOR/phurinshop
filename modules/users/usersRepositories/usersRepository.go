package usersRepositories

import (
	"github.com/PHURINTOR/phurinshop/modules/users"
	"github.com/PHURINTOR/phurinshop/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

// -------------------------------------- Interface ---------------------------------
type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
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

// ----------------------------------------------- Insert
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
