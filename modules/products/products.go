package products

import (
	"github.com/PHURINTOR/phurinshop/modules/appinfo"
	"github.com/PHURINTOR/phurinshop/modules/entities"
)

type Product struct {
	Id          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Category    *appinfo.Category  `json:"category"`
	CreatedAt   string             `json:"created_at"` // ดึงเป็น time.time ก็ได้
	UpdatedAt   string             `json:"updated_at"`
	Price       float64            `json:"price"`
	Images      []*entities.Images `json:"images"`
}

// -------------------- Find Product (Array)
type ProductFilter struct {
	Id                      string `query:"id"`     //param
	Search                  string `query:"search"` // title & description
	*entities.PaginationReq        //ประกาศแบบนี้ไม่ต้องใส่ตัวแปร เราไม่ต้อง .ตัวแปรหลายครั้ง สามารถเข้าใช้ stuct ได้เลย
	*entities.SortReq
}
