package appinfoRepositories

import "github.com/jmoiron/sqlx"

// -------------------------------------- Interface ---------------------------------
type IAppinfoRepository interface {
}

//-------------------------------------- Struct--------------------------------
type appinfoRepository struct {
	db *sqlx.DB
}

// -------------------------------------- Constructor--------------------------------
func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{
		db: db,
	}
}

//===========================================================================================
