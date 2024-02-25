package logs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PHURINTOR/phurinshop/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// ------------------------------- Interface ----------------------------------------
type IPhurinLogger interface {
	Print() IPhurinLogger
	Save()
	SetQuery(c *fiber.Ctx) //เอาไว้ประกอบตอนค้นหาส่วนประกอบใน log
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

// ------------------------------- Struct (encap)----------------------------------------
type phurinlogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"statuscode"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

//Inital Constructor

func InitPhurinlogger(c *fiber.Ctx, res any) IPhurinLogger {
	log := &phurinlogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(), //บางครั้ง อาจไม่แสดง ip ขึ้นอยู่กับ revert proxy
		Method:     c.Method(),
		StatusCode: c.Response().StatusCode(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

//------------------------------------------------------- implement Function missing

func (l *phurinlogger) Print() IPhurinLogger {
	//สร้าง function debug json string ใน Util ก่อน
	utils.Debug(l)
	return l
}

// ------------------------------------------------------- implement Function missing
func (l *phurinlogger) Save() {
	//สร้าง function Output data to byte ใน Util ก่อน
	data := utils.Output(l)

	// จะเป็นการเขียนไฟล์  สร้าง asset/log ก่อน
	filename := fmt.Sprintf("./asset/log/phurinlogger_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	//"-" ตัวที่สนใจจะ replace,  ผลลัพธ์เมื่อ replace ไปแล้ว

	//Open file
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file : %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n") //แปลงจาก byte เป็น string ก่อนเขียนลงในไฟล์

}

func (l *phurinlogger) SetQuery(c *fiber.Ctx) {
	var body any

	if err := c.QueryParser(&body); err != nil { //
		log.Printf("Query parser error : %v", err)
	}
	l.Query = body // อะโลเขตลง body แล้ว

}

func (l *phurinlogger) SetBody(c *fiber.Ctx) {
	var body any
	//เวลาที่รับ body ่ผ่าน fiber เข้ามาจะใช้ body parser
	if err := c.BodyParser(&body); err != nil { //
		log.Printf("body Parser error : %v", err)
	}

	//จะสามารถ Setbody ได้ Method ที่ส่งมาจะต้องเป็น PATCH, PUT, POST
	//Get, Delete จะไม่มี
	//ควรดัก path ด้วย เพราะหากดัก password ลง log คงไม่เหมาะสม

	switch l.Path {
	case "v1/users/signup":
		l.Body = "never gonna give you up"
	default:
		l.Body = body //setแล้ว
	}
}

func (l *phurinlogger) SetResponse(res any) {
	l.Response = res
}

//นำไปแทรกใช้กับ Response ใน Success, Error
