package appinfoRepositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/PHURINTOR/phurinshop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

// ======================================= Interface =========================================
type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) // Find Category
	InsertCategory(req []*appinfo.Category) error                          // Insert Category
	DeleteCategory(categoryId int) error                                   // Delete Category
}

// ======================================= Struct ============================================
type appinfoRepository struct {
	db *sqlx.DB
}

// ======================================= Constructor =======================================
func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{
		db: db,
	}
}

// =======================================Missing Function =======================================

// -------------------------------------- Find Category -------------------------------------
func (r *appinfoRepository) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "categories"
	`
	//Concat String
	filterValues := make([]any, 0) //เพราะ selectรับเป็น interface เลยสร้างแบบ any ดีกว่า
	if req.Title != "" {
		query += `
		WHERE (LOWER("title") LIKE $1)`

		filterValues = append(filterValues, "%"+strings.ToLower(req.Title)+"%") //เพิ่มเข้าไปด้านหน้า
	}
	query += ";"

	category := make([]*appinfo.Category, 0)
	if err := r.db.Select(&category, query, filterValues...); err != nil {
		return nil, fmt.Errorf("select categories failed: %v", err)
	}
	return category, nil

}

// -------------------------------------- Insert Category -------------------------------------
func (r *appinfoRepository) InsertCategory(req []*appinfo.Category) error { //array pointer จะส่งข้อมูลหากมี error กลับมาจะได้สะดวกเพราะว่าเป็น Pointer ของตัวแปร
	/*Note :  *[]  = pointer to array    pointer ตัวหนึ่ง *  ชี้ไปหา Array ก้อนหนึ่ง
	  []* = array of pointer    array ก้อนหนึ่ง เก็บ pointer ไว้หลายๆ ตัว
	*/

	//**************** ที่ต้อง Insert หลายๆ ค่าเพราะ รับตัวแปรมาจาก URL หลายๆ ตัว

	ctx := context.Background()
	query := `INSERT INTO "categories" ("title") VALUES`
	//insert แบบ transection คือถ้าไม่สำเร็จ จะ Rollbackใหม่
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	valuesStack := make([]any, 0) //ตัวแปรที่เป็น  Array

	for i, cat := range req { //ต้อง insert ตาม Value string ข้างบน  ตามจำนวน | SQL String ===> valuesStack(any array)
		valuesStack = append(valuesStack, cat.Title)

		if i != len(req)-1 {
			query += fmt.Sprintf(`($%d),`, i+1) // multi value sql String
		} else {
			//query += fmt.Sprintf(`($%d);`, i+1) // ;last value sql String

			query += fmt.Sprintf(`($%d)`, i+1) // อยากให้ insert แล้ว return ID
		}
	}

	query += `
	RETURNING "id";`

	//------------------- Insert Part
	/* Note: =====
	QueryRowxContext(ctx, query, valuesStack...).StructScan()   = SCAN Stuct ได้
	QueryRowContext(ctx, query, valuesStack...).Scan()			= SCAN String ธรรมดา
	*/
	rows, err := tx.QueryxContext(ctx, query, valuesStack...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert categories failed: %v", err)
	}

	var i int
	for rows.Next() { //Disign pattern loop
		if err := rows.Scan(&req[i].Id); err != nil {
			return fmt.Errorf("scan categories id failed: %v", err)
		}
		i++
	}

	//insert ok = Commit
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// -------------------------------------- Delete Category -------------------------------------
func (r *appinfoRepository) DeleteCategory(categoryId int) error {
	ctx := context.Background()
	query := `DELETE FROM "categories" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(ctx, query, categoryId); err != nil {
		return fmt.Errorf("delete categories failed: %v", err)
	}
	return nil
}
