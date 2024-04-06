package orders

import "github.com/PHURINTOR/phurinshop/modules/products"

type Oders struct {
	Id          string           `db:"id" json:"id"`
	UserId      string           `db:"id" json:"user_id"`
	TranferSlip *TranferSlip     `db:"tranfer_slip" json:"tranfer_slip"`
	Products    []*ProductsOrder `json:"products"` // Product Orders
	Address     string           `db:"address" json:"address"`
	Contact     string           `db:"contact" json:"contact"`
	Status      string           `db:"status" json:"status"`
	TotalPaid   float64          `db:"total_paid" json:"total_paid"`
	CreatedAt   string           `db:"created_at" json:"created_at"`
	UpdatedAt   string           `db:"updated_at" json:"updated_at"`
}

type TranferSlip struct {
	Id        string `json:"id"`
	FileName  string `json:"filename"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductsOrder struct {
	Id      string            `db:"id" json:"id"`
	Qty     int               `db:"qty" json:"qty"`
	Product *products.Product `db:"product" json:"product"`
}
