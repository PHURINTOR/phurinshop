package productPatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/PHURINTOR/phurinshop/modules/products"
	"github.com/jmoiron/sqlx"
)

// ใช้ Builder Pattern
// =================================================== Builder ======================================
// ========================= Interface =======================
type IInsertProductsBuilder interface {
	// Input
	initTransaction() error
	insertProduct() error
	insertCategory() error
	insertAttachment() error
	commit() error

	// Output
	getProductId() string
}

// ========================= Struct Builder==============
type insertProductsBuilder struct {
	db  *sqlx.DB          // Database
	tx  *sqlx.Tx          // เอาไว้ใช้ soft delete
	req *products.Product // Input
}

// ========================= Constructor Builder==============
func InsertProductBuilder(db *sqlx.DB, req *products.Product) IInsertProductsBuilder { // รับสองตัวเพราะเดี๋ยว Tx สร้างเอา
	return &insertProductsBuilder{
		db:  db,
		req: req,
	}
}

// ========================= Missing Func Builder==============
// ------------ init Transaction ------
func (b *insertProductsBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

// ------------ init Transaction ------
func (b *insertProductsBuilder) insertProduct() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "products" (
		"title",
		"description",
		"price"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := b.tx.QueryRowContext( //QueryRowContext สามารถ scan ต่อได้เลย,  QueryRow ต้องไป for Row . Next และไม่ต้อง close.db ด้วย
		ctx, //QueryRowxContext สามารถ scan แบบ struct ได้
		query,
		b.req.Title,
		b.req.Description,
		b.req.Price,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product failed: %v", err)
	}
	return nil
}

// ------------ Insert Category  ------
// ไม่ใช่การสร้าง Category ใหม่แต่เป็นการระบุว่า product อยู่ใน category ใด
func (b *insertProductsBuilder) insertCategory() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "products_categories"(
		"product_id",
		"category_id"
	)
	VALUES($1, $2);`

	//ไม่ต้องการรีเทรนใช้ execContext ได้เลย
	if _, err := b.tx.ExecContext(
		ctx,
		query,
		b.req.Id,
		b.req.Category.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product category Failed: %v", err)
	}
	return nil
}

// ------------------------- Insert Images --------
//
//	call --> Upload api = linkPic ---> insertAttachment
func (b *insertProductsBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "images" (
		"filename",
		"url",
		"product_id"
	)
	VALUES`

	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Images {
		valueStack = append(valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

// Commit TX of DB
func (b *insertProductsBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Get Prodcut
func (b *insertProductsBuilder) getProductId() string {
	return b.req.Id
}

// =================================================== Engineer ======================================
type insertProductsEngineer struct {
	builder IInsertProductsBuilder
}

func InsertProductsEngineer(builder IInsertProductsBuilder) *insertProductsEngineer {
	return &insertProductsEngineer{builder: builder}
}

func (en *insertProductsEngineer) InsertProduct() (string, error) {
	// Init Transaction
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	// Insert Product
	if err := en.builder.insertProduct(); err != nil {
		return "", err
	}
	// Insert Category
	if err := en.builder.insertCategory(); err != nil {
		return "", err
	}
	// Insert Attachment
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	// Commit
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getProductId(), nil
}
