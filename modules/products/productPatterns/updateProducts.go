package productPatterns

import (
	"context"
	"fmt"

	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/files"
	"github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	"github.com/PHURINTOR/phurinshop/modules/products"
	"github.com/jmoiron/sqlx"
)

// =================================================== Builder ======================================
// ============ Builder Interface =====
type IUpdateProductsBuilder interface {
	initTransaction() error

	initQuery()
	closeQuery() //concat str									ng

	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategoryQuery() error //คนละ Talble update error แยกไปเลย

	// Images
	insertImages() error
	getOldImages() []*entities.Images // เป็น pointer ด้วยกรณีรองรับเป็น null
	deleteOldImages() error

	updateProducts() error
	getQueryFields() []string // Get ว่าฟิวด์ไหนถูก update บ้าง
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

// ============ Builder Struct =====
type updateProductsBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *products.Product
	filesUsecases  filesUsecases.IFilesUsecase //ลบเพิ่มไฟล์
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

// ============ Builder Constructor =====

func UpdateProdcutsBuilder(db *sqlx.DB, req *products.Product, filesUseses filesUsecases.IFilesUsecase) IUpdateProductsBuilder {
	return &updateProductsBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUseses,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}

}

// ================== Builder missing Func  ====
// -------- Init Tx --------------
func (b *updateProductsBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

// -------- Init query --------------
func (b *updateProductsBuilder) initQuery() {
	b.query += `
	UPDATE "products" SET`
}

// -------- Close query --------------
func (b *updateProductsBuilder) closeQuery() {

	b.values = append(b.values, b.req.Id) // Add data
	b.lastStackIndex = len(b.values)      // นับเพื่อใช้เป็นลำดับ อากิวเม้นต์ $		 update len

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex) // =====> $(len) = $1, $2 ไปเรื่อยๆ

} //concat string

// --------updateTitleQuery --------------
func (b *updateProductsBuilder) updateTitleQuery() {
	// จะ  Stack query String (initQuery) + values: [] any ที่ส่งเข้ามา
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title) // Add data
		b.lastStackIndex = len(b.values)         // นับเพื่อใช้เป็นลำดับ อากิวเม้นต์ $		 update len

		b.queryFields = append(b.queryFields, fmt.Sprintf(`	
		"title" = $%d`, b.lastStackIndex)) // =====> $(len) = $1, $2 ไปเรื่อยๆ
	}
}

// --------updateDescriptionQuery --------------
func (b *updateProductsBuilder) updateDescriptionQuery() {
	// จะ  Stack query String (initQuery) + values: [] any ที่ส่งเข้ามา
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description) // Add data
		b.lastStackIndex = len(b.values)               // นับเพื่อใช้เป็นลำดับ อากิวเม้นต์ $		 update len

		b.queryFields = append(b.queryFields, fmt.Sprintf(`	
			"description" = $%d`, b.lastStackIndex)) // =====> $(len) = $1, $2 ไปเรื่อยๆ
	}
}

// --------updatePriceQuery --------------
func (b *updateProductsBuilder) updatePriceQuery() {
	// จะ  Stack query String (initQuery) + values: [] any ที่ส่งเข้ามา
	if b.req.Price != 0 && b.req.Price < 0 {
		b.values = append(b.values, b.req.Price) // Add data
		b.lastStackIndex = len(b.values)         // นับเพื่อใช้เป็นลำดับ อากิวเม้นต์ $		 update len

		b.queryFields = append(b.queryFields, fmt.Sprintf(`	
			"price" = $%d`, b.lastStackIndex)) // =====> $(len) = $1, $2 ไปเรื่อยๆ
	}
}

// --------updateCategoryQuery --------------
func (b *updateProductsBuilder) updateCategoryQuery() error { //คนละ Talble update error{return nil} แยกไปเลย
	if b.req.Category == nil {
		return nil
	}

	if b.req.Category.Id == 0 {
		return nil
	}

	query := `
	UPDATE "products_categories" SET
		"category_id" = $1
	WHERE "product_id" = $2;`
	//เนื่องจากไม่ได้ขึ้นกับอันเดิม  Update แยกส่วนทำให้รู้ ตัวแปรที่แน่นอน ไม่ต้องไล่ลำดับเหมือน len

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Category.Id,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update products_categories failed: %v", err)
	}
	return nil
}

// --------insertImages --------------
func (b *updateProductsBuilder) insertImages() error {
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
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

// --------getOldImages --------------
func (b *updateProductsBuilder) getOldImages() []*entities.Images { // เป็น pointer ด้วยกรณีรองรับเป็น null
	query := `
		SELECT
			"id",
			"filename",
			"url"
		FROM "images"
		WHERE "product_id" = $1;`

	images := make([]*entities.Images, 0) // Create const
	if err := b.db.Select(
		&images,
		query,
		b.req.Id,
	); err != nil {
		return make([]*entities.Images, 0)
	}
	return images
}

// --------deleteOldImages --------------
func (b *updateProductsBuilder) deleteOldImages() error {
	// import fileusecase file storage  + func get old images
	query := `DELETE FROM "images" WHERE "product_id" = $1;`

	//Check old images  loop delete
	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/products/%s", img.FileName),
			})
		}
		//ลบไปแล้วไม่ต้อง ดัก err
		// if err := b.filesUsecases.DeleteFileGCP(deleteFileReq); err != nil {
		// 	b.tx.Rollback()
		// 	return err
		// }
		b.filesUsecases.DeleteFileGCP(deleteFileReq)
	}

	// use
	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete images failed: %v", err)
	}
	return nil
}

// --------updateProducts --------------
func (b *updateProductsBuilder) updateProducts() error {
	// เป็นเพียง func ว่าอนุญาติรัน Update ได้ไหม ส่วน sql query จะไปทำ fuction อื่น
	fmt.Printf(b.query)
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update product failed: %v", err)
	}
	return nil
}

// --------getQueryFields --------------
func (b *updateProductsBuilder) getQueryFields() []string { // Get List ว่าฟิวด์ไหนถูก update บ้าง

	return b.queryFields
}

// --------getValues --------------
func (b *updateProductsBuilder) getValues() []any { return b.values }

// --------getQuery --------------
func (b *updateProductsBuilder) getQuery() string { return b.query }

// --------setQuery --------------
func (b *updateProductsBuilder) setQuery(query string) { b.query = query }

// --------getImagesLen --------------
func (b *updateProductsBuilder) getImagesLen() int { return len(b.req.Images) }

// -------- Commit --------------
func (b *updateProductsBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

// =================================================== Engineer ======================================
// ================ Engineer struct
type updateProductsEngineer struct {
	builder IUpdateProductsBuilder
}

// ================ Engineer Constructor
func UpdateProductEngineer(b IUpdateProductsBuilder) *updateProductsEngineer {
	return &updateProductsEngineer{builder: b}
}

// ================ Engineer missing Func summary sql
func (en *updateProductsEngineer) sumQueryFields() {
	en.builder.updateTitleQuery()
	en.builder.updateDescriptionQuery()
	en.builder.updatePriceQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery() //ของเดิม
		if i != len(fields)-1 {        //เช็คว่าเป็นรอบสุดท้ายหรือไม่
			en.builder.setQuery(query + fields[i] + ",") // query ole + new query
		} else {
			en.builder.setQuery(query + fields[i]) // กรณีเป็นตัวสุดท้าย
		}
	}

}

// ================ Engineer missing Func Update Products
func (en *updateProductsEngineer) UpdateProduct() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	// Update Product
	if err := en.builder.updateProducts(); err != nil {
		return err
	}

	// Update category
	if err := en.builder.updateCategoryQuery(); err != nil {
		return err
	}

	if en.builder.getImagesLen() > 0 {
		// delete old images
		if err := en.builder.deleteOldImages(); err != nil {
			return err
		}
		// insert new images
		if err := en.builder.insertImages(); err != nil {
			return err
		}
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}

	return nil
}
