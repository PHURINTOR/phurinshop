package middlewaresRepositories

import "github.com/jmoiron/sqlx"

//-------------------------------------------- Interface --------------------------
type IMiidlewareRepository interface {
}

//------------------------------------ Struct ------------------------------------
type middlewareRepository struct {
	db *sqlx.DB
}

//------------------------------- Constructor -------------------------------------
func MiddlewareRepository(db *sqlx.DB) IMiidlewareRepository {
	return &middlewareRepository{
		db: db,
	}
}
