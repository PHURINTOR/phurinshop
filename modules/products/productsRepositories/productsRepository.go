package productsRepositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/files/filesUsecases"
	"github.com/PHURINTOR/phurinshop/modules/products"
	"github.com/PHURINTOR/phurinshop/modules/products/productPatterns"
	"github.com/jmoiron/sqlx"
)

// ======================================= Interface =========================================
type IProductsRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) ([]*products.Product, int)
	InsertProducts(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
}

// ======================================= Struct ============================================
type productsRepository struct {
	db            *sqlx.DB
	cfg           config.IConfig
	filesUsecases filesUsecases.IFilesUsecase
}

// ======================================= Constructor =======================================
func ProductsRepository(db *sqlx.DB, cfg config.IConfig, fileUsecases filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:            db,
		cfg:           cfg,
		filesUsecases: fileUsecases,
	}
}

// ======================================= Missing Function =======================================
// ------------  FindOneProduct ---------------

// ** *COALESCE เช็คค่าว่าง
func (r *productsRepository) FindOneProduct(productId string) (*products.Product, error) {
	fmt.Println(productId)
	query := `
	SELECT
		to_jsonb("t")
	FROM (
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
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";`

	// inital const
	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Images, 0),
	}

	// query product
	if err := r.db.Get(&productBytes, query, productId); err != nil { //retrun bytes type
		return nil, fmt.Errorf("get Product failed: %v", err)
	}

	// Covert Type query
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshl product failed :%v", err)
	}
	return product, nil
}

// ------------------------------- Find Product

func (r *productsRepository) FindProducts(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productPatterns.FindProductBuilder(r.db, req)
	engineer := productPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	return result, count
}

// ------------------------------- Insert Product
func (r *productsRepository) InsertProducts(req *products.Product) (*products.Product, error) {

	builder := productPatterns.InsertProductBuilder(r.db, req)
	productId, err := productPatterns.InsertProductsEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}
	// insert แล้วจะได้ Return ID กลับมา  ให้เอาไปใช้ findOne เป็นการรียูส func

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// ------------------------------- Update Product
func (r *productsRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productPatterns.UpdateProdcutsBuilder(r.db, req, r.filesUsecases)
	engineer := productPatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	// ok
	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// ------------------------------- Delete Product
func (r *productsRepository) DeleteProduct(productId string) error {
	// Delete 2 ส่วน 1. Database local 2. GCP buket (Handler)

	// 1. local delete  ----> Repository
	//
	query := `DELETE FROM "products" WHERE "id" = $1;`
	if _, err := r.db.ExecContext(context.Background(), query, productId); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}
	return nil
}
