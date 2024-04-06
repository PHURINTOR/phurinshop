package productPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PHURINTOR/phurinshop/modules/products"
	"github.com/PHURINTOR/phurinshop/pkg/utils"
	"github.com/jmoiron/sqlx"
)

// Builder Patterns

// interface Builder, type Builder, constructor Builder  ,  Engineer
// Builder --> Engineer  ---> Create
// =================================================== Builder ======================================
// ---------------- Builder Interface -------------
type IFindProductBuilder interface {
	openJsonQuery()
	initQuery()  // ดึงข้อมูลปกติ
	countQuery() // นับข้อมูลเฉยๆ
	whereQuery() // คิวรีเงื่อนไข
	sort()       // เรียง
	paginate()
	closeJsonQuery()
	resetQuery()                 // รีค่าใน stuct ค้าง
	Result() []*products.Product //ผลลัพธ์
	Count() int
	PrintQuery()
}

// ---------------- Builder Stuct -------------
type findProductBuilder struct {
	db             *sqlx.DB
	req            *products.ProductFilter
	query          string //อาจมีการคิวรีที่ซับซ้อน  ทำให้ประกาศในนี้ไว้ประกอบ
	lastStackIndex int
	values         []any // อากิวเม้นหลายๆ ตัว กัน sql injection
}

// ---------------- Builder Constructor -------------
func FindProductBuilder(db *sqlx.DB, req *products.ProductFilter) IFindProductBuilder {
	return &findProductBuilder{
		db:  db,
		req: req,
	}
}

// ---------------- Builder Missing Function -------------
// เหมือนแยกชุดคำสั่งไว้ แล้วค่อยเอามาประกอบกันเป็นคิวรี่

// --------  jsonOpen
func (b *findProductBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

// -------- jsonClose
func (b *findProductBuilder) closeJsonQuery() {
	b.query += `
) AS "t";`
}

// ------------- Sql Find Product
func (b *findProductBuilder) initQuery() { // ดึงข้อมูลปกติ
	b.query += `
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			(
				SELECT
					to_jsonb("ct")
				FROM (
					SELECT
						"c"."id",
						"c"."title"
					FROM "categories" "c"
						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "products" "p"
		WHERE 1 = 1`
}

// ------- Count
func (b *findProductBuilder) countQuery() { // นับข้อมูลเฉยๆ ต้อง query มาก่อน |  query ---> reset query ---> CountQuery
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "products" "p"
		WHERE 1=1`

	// 1=1 คือ ให้เท่ากับ True เนื่องจากเราจะใช้ parameter filter ตัวเดียวกันกับที่ส่งมาร่วมกับ query อื่น
}

// ------- คิวรีเงื่อนไข
func (b *findProductBuilder) whereQuery() {
	// reset qeuey and value
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// ID Check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)
		queryWhereStack = append(queryWhereStack, `
		AND "p"."id" = ?`)
	}

	// Search Check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)
		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("p"."title") LIKE ? OR LOWER("p"."description") LIKE ?)`) // รอ convert string ข้างล่าง
	}

	// Concat string
	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1) //convert to $1
		} else {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
			queryWhere = strings.Replace(queryWhere, "?", "$"+strconv.Itoa(i+2), 1)
		}

		// Last Stack Record
		b.lastStackIndex = len(b.values)

		// Summary Qeury
		b.query += queryWhere
	}

}

// -------เรียง
func (b *findProductBuilder) sort() {
	orderByMap := map[string]string{
		"id":    "\"p\".\"id\"",
		"title": "\"p\".\"title\"",
		"price": "\"p\".\"price\"",
	}
	// ถ้า ไม่ได้ใส่ type sort มาก็จะเป็นค่าเริ่มต้น
	if orderByMap[b.req.OrderBy] == "" {
		b.req.OrderBy = orderByMap["title"]
	} else {
		b.req.OrderBy = orderByMap[b.req.OrderBy]
	}

	//มากไปน้อย น้อยไปมาก
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[b.req.Sort] == "" {
		b.req.Sort = sortMap["ASC"]
	} else {
		b.req.Sort = sortMap[strings.ToUpper(b.req.OrderBy)]
	}
	b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`
		ORDER BY $%d %s`, b.lastStackIndex+1, b.req.Sort)
	b.lastStackIndex = len(b.values)
}
func (b *findProductBuilder) paginate() {
	// offset (page -1)*limit

	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)

}

// ------- ResetQuery
// รีค่าใน stuct ค้าง
func (b *findProductBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
}

// ------- Result-----------------------------
func (b *findProductBuilder) Result() []*products.Product {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// sql command ---> sqlx = bytes ----> convert to json
	bytes := make([]byte, 0)
	productsData := make([]*products.Product, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find product failed: %v\n", err)
		return make([]*products.Product, 0)
	}

	// Result convert unmarshal json
	if err := json.Unmarshal(bytes, &productsData); err != nil {
		log.Printf("unmarshal products failed: %v", err)
		return make([]*products.Product, 0)
	}
	b.resetQuery()
	return productsData
}

// ------- Count Result-----------------------------
func (b *findProductBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	//
	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count product failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}
func (b *findProductBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

// =================================================== Engineer ======================================

// ----- Type
type findProductEngineer struct {
	builder IFindProductBuilder
}

// ----- Constructor ---
func FindProductEngineer(builder IFindProductBuilder) *findProductEngineer {
	return &findProductEngineer{builder: builder}
}

// ---- Use ----
func (en *findProductEngineer) FindProduct() IFindProductBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findProductEngineer) CountProduct() IFindProductBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
