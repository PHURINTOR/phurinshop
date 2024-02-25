package entities

import (
	"github.com/PHURINTOR/phurinshop/pkg/logs"
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, tractId string, msg string) IResponse
	Res() error
}

// ---------------------------------------------------------------------- Struct Return Response -----------------------------------
type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool //เอาไว้เช็คว่าสำเร็จหรือไม่

}

// ----------------------------------------------------------------------
type ErrorResponse struct {
	TraceId string `json:"trace_id"` //id ของ error นั้นๆ
	Msg     string `json:"message"`
}

// ---------------------------------------------------------------------- Constructor  Response -------------------------------------
func NewErrorResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c, //ครอบฟังก์ชัน ของ package คือ เรียก context จาก fiber.ctx มาครอบเพื่อปรับแต่ง
	}
}

// ----------------------------------------------------------------------  Implement function --------------------------------------
func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	logs.InitPhurinlogger(r.Context, &r.Data).Print().Save() //เพิ่มมาหลังจากสร้าง logs ของตัวเอง
	return r
}

func (r *Response) Error(code int, tractId string, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: tractId,
		Msg:     msg,
	}
	r.IsError = true
	logs.InitPhurinlogger(r.Context, &r.ErrorRes).Print().Save() //เพิ่มมาหลังจากสร้าง logs ของตัวเอง
	return r
}

func (r *Response) Res() error {
	return r.Context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorRes
		}
		return &r.Data
	}())
}

//Use แบบเดิม
// c.Status(fiber.StatusOK).JSON(res)

//แบบใหม่ ทำเพื่อให้สามารถเก็บ log ได้ง่ายยิ่งขึ้น
//	entiites.NewResponse(c).Success(fiber.StatusOK, res).Res()
