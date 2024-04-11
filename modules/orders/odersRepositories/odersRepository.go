package odersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PHURINTOR/phurinshop/modules/orders"
	"github.com/PHURINTOR/phurinshop/modules/orders/ordersPatterns"
	"github.com/jmoiron/sqlx"
)

// ======================================= Interface =========================================
type IOrdersRepository interface {
	FindOneOrder(orderId string) (*orders.Oders, error)
	FindOrder(req *orders.OrderFilter) ([]*orders.Oders, int)
	InsertOrder(req *orders.Oders) (string, error)
	UpdateOrder(req *orders.Oders) error
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

// ----------------------- FindOneOrder -------------------
func (r *ordersRepository) FindOneOrder(orderId string) (*orders.Oders, error) {
	// (json array)
	// การเก็บข้อมูล Product ใน Order จะเก็บทั้งหมด  จะไม่ใช้การเชื่อมโยง เพราะเวลาแก้ไข Product --> Order ที่ทำการไปแล้วจะไม่กระทบเปลี่ยนแปลง
	// order(id) อยู่ใน product_orders (1:M)  --- >
	// product ใน product_orders = json ทั้งก้อน

	/*		SELECT
			SUM("po"."product"->>'price')  |   ->> เข้าถึงฟิวด์ json -> เข้าถึง index ด้วย sql
			SUM(("po"."product"->>'price')::FLOAT)   = ::FLOAT คือแปลงค่าเนื่องจาก ๋json เป็น text

			SELECT
					SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))

					COALESCE เอาไว้ตรวจค่า Null ใน sql ถ้า null จะให้เท่ากับ 0
	*/
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"o"."id",
			"o"."user_id",
			"o"."transfer_slip",
			"o"."status",
			(
				SELECT
					array_to_json(array_agg("pt"))
				FROM (
					SELECT
						"spo"."id",
						"spo","qty",
						"spo","product"
					FROM "products_orders" "spo"
					WHERE "spo"."order_id" = "o"."id"
				) AS "pt"
			) AS "products",
			"o"."address",
			"o"."contact",
			(
				SELECT
					SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))
				FROM "products_orders" "po"
				WHERE "po"."order_id" = "o"."id"
			) AS "total_paid",
			"o"."created_at",
			"o"."updated_at"
		FROM "orders" "o"
		WHERE "o"."id" = $1
	) AS "t";`

	// inital
	orderData := &orders.Oders{ //ประกาศ type ที่มีการอ้างอิงตำแหน่งก่อน
		TranferSlip: &orders.TranferSlip{},
		Products:    make([]*orders.ProductsOrder, 0),
	}

	// Query
	raw := make([]byte, 0)
	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("get order fialed: %v", err)
	}

	// output Convert(byte --> Unmarshal --> json)
	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return orderData, nil
}

// ----------------------- FindManyOrder -------------------

func (r *ordersRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Oders, int) {
	builder := ordersPatterns.FindOrderBuilder(r.db, req)
	engineer := ordersPatterns.FindOrderEngineer(builder)
	return engineer.FindOrder(), engineer.CountOrder()
}

// ----------------------- Insert Orders -------------------
func (r *ordersRepository) InsertOrder(req *orders.Oders) (string, error) {
	builder := ordersPatterns.InsertOrderBuilder(r.db, req)
	orderId, err := ordersPatterns.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}
	return orderId, nil
}

// ----------------------- Update Orders -------------------
func (r *ordersRepository) UpdateOrder(req *orders.Oders) error {
	query := `
	UPDATE "orders" SET`

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	// 1. Update Status
	if req.Status != "" {
		values = append(values, req.Status)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"status" = $%d?`, lastIndex))
		lastIndex++
	}

	// 2. Update Tranfer slip
	if req.TranferSlip != nil {
		values = append(values, req.TranferSlip)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"transfer_slip" = $%d?`, lastIndex))
		lastIndex++
	}

	// Where = id
	values = append(values, req.Id)

	queryClose := fmt.Sprintf(`
	WHERE "id" = $%d;`, lastIndex)

	// Loop query Stack check last
	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1) //ต่อ
		}
	}

	// Summary Query
	query += queryClose

	fmt.Println(query)

	// Update excute
	if _, err := r.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("update error failed: %v", err)
	}
	return nil
}
