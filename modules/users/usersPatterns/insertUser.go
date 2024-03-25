package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PHURINTOR/phurinshop/modules/users"
	"github.com/jmoiron/sqlx"
)

//Method Factory

// 1. โรงงานใหญ่ก่อน
type IInsertUser interface {
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	//เลือกอันใดอันหนึ่ง  แล้ว Return + result
	Result() (*users.UserPassport, error)
}

// -------------------------------- Struct --------------------
type userReq struct {
	id  string //เมื่อ register แล้ว id ที่เป็น autorun จะreturn กลับมา
	req *users.UserRegisterReq
	db  *sqlx.DB
}

// -------------------------------- Struct  Factory --------------------
type customer struct {
	*userReq //ไม่มีตัวแปรข้างหน้าคือ Stuct นี้สามารถเข้าถึง id, req, db ได้โดยตรงเลย
}

type admin struct {
	*userReq //ไม่มีตัวแปรข้างหน้าคือ Stuct นี้สามารถเข้าถึง id, req, db ได้โดยตรงเลย
}

// =========================================Constuctor===============================================
func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) IInsertUser { //factory ใหญ๋
	//isAdmin อาจใช้ Enum ได้กรณีมีหลายตัว

	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser { //factory ย่อย
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser { //factory ย่อย
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

// ========================== implement missing function ==============
func (f *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //ถ้าไม่เสร็จภายใน 5 วิ จะยกเลิกและคืนทรัพยากร
	defer cancel()

	query := `
	INSERT INTO "users"(
		"email",
		"password",
		"username",
		"role_id"
	)VALUES($1, $2, $3, 1)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil { //scan คือ สิ่งที่อยากรีเทรินออกมา = id
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return f, nil
}

func (f *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //ถ้าไม่เสร็จภายใน 5 วิ จะยกเลิกและคืนทรัพยากร
	defer cancel()
	query := `
	INSERT INTO "users"(
		"email",
		"password",
		"username",
		"role_id"
	)VALUES($1, $2, $3, 2)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil { //scan คือ สิ่งที่อยากรีเทรินออกมา = id
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return f, nil
}

func (f *userReq) Result() (*users.UserPassport, error) {
	//Get data
	query := `
	SELECT
		json_build_object(
			'user',"t",
			'token', NULL
		)
	FROM (
		SELECT
			"u"."id",
			"u"."email",
			"u"."username",
			"u"."role_id"
		FROM	"users" "u"
		WHERE "u"."id"= $1
	)AS "t"`

	data := make([]byte, 0)
	if err := f.db.Get(&data, query, f.id); err != nil { //.Get สามารถนำเข้าแบบ interface ได้เลยไม่จำเป็นต้องใช้ select
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil { //Unmarshal มีคุณสมบัติแปลง oject ใดๆ ก็ตาม struct ที่เรากำหนด
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}
	return user, nil
}

//**** นำเอาไปใช้ใน Repository ใน interface
