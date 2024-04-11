package ordersPatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/PHURINTOR/phurinshop/modules/orders"
	"github.com/jmoiron/sqlx"
)

// =================================================== Builder ======================================
// ---------------- Builder Interface -------------
type IInsertOrderBuilder interface {
	initTransaction() error
	insertOrder() error
	insertProductsOrder() error
	getOrderId() string
	commit() error
}

// ---------------- Builder Stuct ------------------------
type insertOrderBuilder struct {
	req *orders.Oders
	db  *sqlx.DB
	tx  *sqlx.Tx
}

// ---------------- Builder Constructor ------------------
func InsertOrderBuilder(db *sqlx.DB, req *orders.Oders) IInsertOrderBuilder {
	return &insertOrderBuilder{
		db:  db,
		req: req,
	}
}

// =================================================== Engineer ======================================
// ---------------- Engineer Stuct ------------------------
type insertOrderEngineer struct {
	builder IInsertOrderBuilder
}

// ---------------- Engineer Constructor ------------------------
func InsertOrderEngineer(b IInsertOrderBuilder) *insertOrderEngineer {
	return &insertOrderEngineer{builder: b}
}

// -------------- Engineer missing Func ----------
func (en *insertOrderEngineer) InsertOrder() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	if err := en.builder.insertOrder(); err != nil {
		return "", err
	}
	if err := en.builder.insertProductsOrder(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getOrderId(), nil
}

// =================================================== Missing Func ======================================
// ============================== Builder =======================
// ---- init -----------
func (b *insertOrderBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

// ---- getOrderID -----------
func (b *insertOrderBuilder) getOrderId() string {

	return ""
}

// ---- insert -----------
func (b *insertOrderBuilder) insertOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	INSERT INTO "orders"(
		"user_id",
		"contact",
		"address",
		"transfer_slip",
		"status"
	)
	VALUES
	($1, $2, $3, $4, $5)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.UserId,
		b.req.Contact,
		b.req.Address,
		b.req.TranferSlip,
		b.req.Status,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert order failed: %v", err)
	}
	return nil
}

// ---- insertProductsOrder -----------
func (b *insertOrderBuilder) insertProductsOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
	INSERT INTO "products_orders"(
		"order_id",
		"qty",
		"product"
	)
	VALUES`

	// ------ Concat String --------------
	value := make([]any, 0)
	lastIndex := 0
	for i := range b.req.Products {
		value = append(
			value,
			b.req.Id,
			b.req.Products[i].Qty,
			b.req.Products[i].Product,
		)
		if i != len(b.req.Products)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, lastIndex+1, lastIndex+2, lastIndex+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, lastIndex+1, lastIndex+2, lastIndex+3)
		}
		lastIndex += 3
	}

	// Excute
	if _, err := b.tx.ExecContext(ctx, query, value...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert porducts_order failed: %v", err)
	}
	return nil
}

// ---- commit -----------
func (b *insertOrderBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}
