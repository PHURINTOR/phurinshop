package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

// มีหน้าที่ในการ Gen Token 1. Access Token, 2. API Token, 3. Admin Token อื่นๆ
// Interface Token
// --------------------------------------- Interface --------------------------
// Main Struct
type IPhurinshopAuth interface {
	SignToken() string //ไม่จำเป็นต้องเรียกใช้งาน Constructor แค่ pass ค่าเข้าไปก็เพียงพอ เรียกดู playload ได้เลย
}

// Interface Factory #1
type IPhurinshopAdmin interface {
	SignToken() string
}

// --------------------------------------- Enum Token --------------------------
type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

// --------------------------------------- Factory Main --------------------------
type phurinshopAuth struct {
	mapClaims *phurinshopMapClaims
	cfg       config.IJwtConfig
}

// --------------------------------------- Factory #1
type phurinshopAdmin struct {
	*phurinshopAuth
}

// --------------------------------------- Factory #2
type phurinshopMapClaims struct {
	Claims               *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims                   //ใน golang-jwt บังคับต้องมี
}

// ================================== Function Plug=============================
// ------ ต้องสร้างเพราะ type ที่เขาต้องการวุ่นวาย ----
func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9))))) //nano sec convert
}

// refresh token วัน Gen token จะเปลียน  แต่วันหมดอายุจะเท่าเดิม  เพราะเราจะต้อง Gen Refresh Token ใหม่เพื่อความปลอดภัย
func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0)) // time ปกติจะมี type = int 64   แต่เรา้องการ time.time  = time.Unix convert sec to time.time
}

// Parse Token
func ParseToken(cfg config.IJwtConfig, tokenString string) (*phurinshopMapClaims, error) {

	//ParseWithClaims มี playload, Parse ไม่มี payload
	token, err := jwt.ParseWithClaims(tokenString, &phurinshopMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})

	//เช็ครูปแบบ
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	//Claim Check
	if claims, ok := token.Claims.(*phurinshopMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
	//Output struct ที่เป็น playload, callback Function เอาไว้ตรวจ อัลกอ ว่า sign มาอย่างไร
}

// Parse Token Admin
func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*phurinshopMapClaims, error) {

	//ParseWithClaims มี playload, Parse ไม่มี payload
	token, err := jwt.ParseWithClaims(tokenString, &phurinshopMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil //***************** สิ่งที่แตกต่าง
	})

	//เช็ครูปแบบ
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	//Claim Check
	if claims, ok := token.Claims.(*phurinshopMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
	//Output struct ที่เป็น playload, callback Function เอาไว้ตรวจ อัลกอ ว่า sign มาอย่างไร
}

// ------------------------RepeatToken
// เราจะรับ Refresh Token เข้ามา  เปลี่ยนวัน Gen ,  วันหมดอายุเท่าเดิม,  เพื่อความปลอดภัย
func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &phurinshopAuth{
		cfg: cfg,
		mapClaims: &phurinshopMapClaims{
			RegisteredClaims: jwt.RegisteredClaims{ //มาจาก jwt-go jwt.RegisteredClaims
				Issuer:    "phurinshopapi",
				Subject:   "access-token",
				Audience:  []string{"costomer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),      //วันหมดอายุเดิมเข้ามา เพื่อ gen ใหม่
				NotBefore: jwt.NewNumericDate(time.Now()), //key นี้จะใช้ได้ตอนไหน
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

// ================================= Function  Main Factory ===============================
// ----------------------- Creator ---------------------------
func NewphurinshopAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IPhurinshopAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("Unknow token type")
	}
}

// ------------------------ ConcretorCreator A1-------------
func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IPhurinshopAuth {
	return &phurinshopAuth{
		cfg: cfg,
		mapClaims: &phurinshopMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{ //มาจาก jwt-go jwt.RegisteredClaims
				Issuer:    "phurinshopapi",
				Subject:   "access-token",
				Audience:  []string{"costomer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()), //key นี้จะใช้ได้ตอนไหน
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

// ------------------------ ConcretorCreator A2-------------
func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IPhurinshopAuth {
	return &phurinshopAuth{
		cfg: cfg,
		mapClaims: &phurinshopMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{ //มาจาก jwt-go jwt.RegisteredClaims
				Issuer:    "phurinshopapi",
				Subject:   "refresh-token",
				Audience:  []string{"costomer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()), //key นี้จะใช้ได้ตอนไหน
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

// ------------------------ ConcretorCreator B-------------

func newAdminToken(cfg config.IJwtConfig) IPhurinshopAuth {
	return &phurinshopAdmin{ //ต้อง return ซ้อนเพราะสืบทอดมา
		phurinshopAuth: &phurinshopAuth{
			cfg: cfg,
			mapClaims: &phurinshopMapClaims{
				Claims: nil, //ไม่จำเป็นต้องมี payload
				RegisteredClaims: jwt.RegisteredClaims{ //มาจาก jwt-go jwt.RegisteredClaims
					Issuer:    "phurinshopapi",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),        // 5min
					NotBefore: jwt.NewNumericDate(time.Now()), //key นี้จะใช้ได้ตอนไหน
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

// ================================= missing Func =============================
func (a *phurinshopAuth) SignToken() string {
	//sign token  ด้วยฟังก์ชั่น missing Jwt ได้เลย แต่ยังไม่ใช่ string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims) //ใช้ key เดียว SigningMethodHS256  , hs Symatic, AES asymatic type
	ss, _ := token.SignedString(a.cfg.SecretKey())                  //อาจใช้ key gen เลขได้
	return ss
}

// ----------------------- Factory #1
func (a *phurinshopAdmin) SignToken() string {
	//sign token  ด้วยฟังก์ชั่น missing Jwt ได้เลย แต่ยังไม่ใช่ string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims) //ใช้ key เดียว SigningMethodHS256  , hs Symatic, AES asymatic type
	ss, _ := token.SignedString(a.cfg.AdminKey())                   //อาจใช้ key gen เลขได้
	return ss

	//เขียน Parse Token ด้วย
}
