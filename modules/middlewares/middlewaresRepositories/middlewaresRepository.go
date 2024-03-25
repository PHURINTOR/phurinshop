package middlewaresRepositories

import (
	"fmt"

	"github.com/PHURINTOR/phurinshop/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

// -------------------------------------------- Interface --------------------------
type IMiidlewareRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error) //Find Role id for Authurize
}

// ------------------------------------ Struct ------------------------------------
type middlewareRepository struct {
	db *sqlx.DB
}

// ------------------------------- Constructor -------------------------------------
func MiddlewareRepository(db *sqlx.DB) IMiidlewareRepository {
	return &middlewareRepository{
		db: db,
	}
}

// ---------------------------- Check User Token Middleware
func (r *middlewareRepository) FindAccessToken(userId, accessToken string) bool {

	/*คำสั่งพิเศษใน Postgre*/
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true
}

// --------------------- Find Role-id -----------------
func (r *middlewareRepository) FindRole() ([]*middlewares.Role, error) {
	query :=
		`
	SELECT
		"id",
		"title"
	FROM "roles"
	ORDER BY "id" DESC;
	`
	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are emtry")
	}
	return roles, nil
}
