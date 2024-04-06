package odersRepositories

import "github.com/jmoiron/sqlx"

// ======================================= Interface =========================================
type IOrdersRepository interface {
}

// ======================================= Struct ============================================
type ordersRepository struct {
	db *sqlx.DB
}

// ======================================= Constructor =======================================
func OdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{db: db}
}

// ======================================= missing Func =======================================
