package odersHandlers

import (
	"strings"
	"time"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/entities"
	"github.com/PHURINTOR/phurinshop/modules/orders"
	"github.com/PHURINTOR/phurinshop/modules/orders/odersUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ======================================= Enum =========================================
type ordersHandlersErrCode string

const (
	findOneOrderErr ordersHandlersErrCode = "orders-001"
	findOrderErr    ordersHandlersErrCode = "orders-002"
	insertOrderErr  ordersHandlersErrCode = "orders-003"
	updateOrderErr  ordersHandlersErrCode = "orders-004"
)

// ======================================= Interface =========================================
type IOdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}

// ======================================= Struct ============================================
type odersHandler struct {
	cfg          config.IConfig
	odersUsecase odersUsecases.IOrdersUsecase
}

// ======================================= Constructor =======================================
func OdersHandler(cfg config.IConfig, odersUsecase odersUsecases.IOrdersUsecase) IOdersHandler {
	return &odersHandler{
		cfg:          cfg,
		odersUsecase: odersUsecase,
	}
}

// ======================================= missing Func =======================================
// ---------- FindOneOrder -------------
func (h *odersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.odersUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewErrorResponse(c).Success(fiber.StatusOK, order).Res()
}

// ---------- FindManyOrders -------------
func (h *odersHandler) FindOrder(c *fiber.Ctx) error {

	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	// varidate defult page value
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Check OrderBy
	orderByMap := map[string]string{
		"id":         `"o","id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	// Sort
	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"start date invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewErrorResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"end date invalid",
			).Res()
		}

		req.EndDate = end.Format("2006-01-02")
	}
	// Usecase
	orders := h.odersUsecase.FindOrder(req)
	return entities.NewErrorResponse(c).Success(fiber.StatusOK, orders).Res()
}

// ---------- Insert Order -------------
func (h *odersHandler) InsertOrder(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string) // เอาไว้ว่า ให้สามารถสร้างได้เฉพาะคนที่เป็นเจ้าของ ID
	req := &orders.Oders{
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	// เช็คว่าไม่ได้ซื้ออะไร
	if len(req.Products) == 0 {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(insertOrderErr),
			"products are emtry",
		).Res()
	}

	// check role
	if c.Locals("userRoleId").(int) != 2 { // ต้องแปลงเป็น Int locals valut = type any   ,  2 =Role admin
		req.UserId = userId
	}

	// check watting

	req.Status = "waiting" // Set status defult value

	req.TotalPaid = 0 // Set Total defult value

	// Execute
	order, err := h.odersUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(insertOrderErr),
			"insert order found",
		).Res()
	}

	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, order).Res()
}

// ---------- Update Order -------------
// admin update status ได้ทั้งหมด , user ได้แค่ status cancel
func (h *odersHandler) UpdateOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	req := new(orders.Oders)

	if err := c.BodyParser(req); err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}

	req.Id = orderId

	// Check status
	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled", // user only
	}

	// Check Role
	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else if strings.ToLower(req.Status) == statusMap["canceled"] { // user
		req.Status = statusMap["canceled"]
	}

	// Update TranferSlip
	if req.TranferSlip != nil {
		if req.TranferSlip.Id == "" {
			req.TranferSlip.Id = uuid.NewString()
		}

		if req.TranferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewErrorResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(updateOrderErr),
					err.Error(),
				).Res()
			}
			now := time.Now().In(loc)
			// YYYY-MM-DD HH:MM:SS
			// YYYY = 2006,   MM = 01,  DD=02   |  HH =15    MM=04  SS=05
			req.TranferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")

		}
	}

	order, err := h.odersUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewErrorResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}
	return entities.NewErrorResponse(c).Success(fiber.StatusCreated, order).Res()
}
