package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// ----------------------------------------- Function Load config (Path) ----------------------------------------------------------------------------------------
func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}
	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			//-------------------------------------------------------------  type Int =-------------------------
			port: func() int {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("load port failed: %v", err)
				}
				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			//-------------------------------------------------------------  type time =-------------------------
			readTimeout: func() time.Duration {
				rt, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("load ReadTimeout failed: %v", err)
				}
				return time.Duration(int64(rt) * int64(math.Pow10(9))) //แปลงค่า เป็น interface แบบเวลา
			}(),

			writeTimeout: func() time.Duration {
				wt, err := strconv.Atoi(envMap["APP_WRITE_TIMEOUT"])
				if err != nil {
					log.Fatalf("load WriteTimeout failed: %v", err)
				}
				return time.Duration(int64(wt) * int64(math.Pow10(9))) //แปลงค่า เป็น interface แบบเวลา
			}(),

			//-------------------------------------------------------------  type Int =-------------------------
			bodyLimit: func() int {
				b, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("load Body Limit failed: %v", err)
				}
				return b
			}(),
			fileLimit: func() int {
				f, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("load File Limit failed: %v", err)
				}
				return f
			}(),
			gcpbucket: envMap["APP_GCP_BUCKET"],
		},

		db: &db{
			host: envMap["DB_HOST"],
			//-------------------------------------------------------------  type Int =-------------------------
			port: func() int {
				dp, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("load DB port failed: %v", err)
				}
				return dp
			}(),
			protocol: envMap["DB_PROTOCOL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			//-------------------------------------------------------------  type Int =-------------------------
			maxConnections: func() int {
				maxCon, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("load MaxConnections failed: %v", err)
				}
				return maxCon
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			secretKey: envMap["JWT_SECRET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpiresAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("load AccessExpiresAt failed: %v", err)
				}
				return t
			}(),
			refreshExpiresAt: func() int {
				rt, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("load RefreshExpiresAt failed: %v", err)
				}
				return rt
			}(),
		},
	}
}

// ---------------------------------------------Config ----------------------------------------------

type IConfig interface {
	App() IAppConfig
	Db() IDBConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

// ------------------------------------------ App -----------------------------------
type IAppConfig interface {
	Url() string //host:port
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	Gcpbucket() string
}
type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration //
	bodyLimit    int           //byte
	fileLimit    int           //byte
	gcpbucket    string
}

func (c *config) App() IAppConfig { //ใช้ Pointer เพราะไม่ต้อง copy เร็วกว่า ดีกว่า copy struct
	return c.app //*** เมื่อ return interface แล้วเรา return object คือ ให้ c.app เข้าถึงข้อมูล ใน interface ได้
}

// implement Functions
func (a *app) Url() string                 { return fmt.Sprintf("%v:%v", a.host, a.port) } //host:port
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) FileLimit() int              { return a.fileLimit }
func (a *app) Gcpbucket() string           { return a.gcpbucket }

// ------------------------------------------ DB  -----------------------------------

type IDBConfig interface { //encapsulation เข้าผ่าน missage
	Url() string //host+port+peotocol + username + password + sslMode
	MaxCon() int
}
type db struct {
	host           string
	port           int
	protocol       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

func (c *config) Db() IDBConfig { //ใช้ Pointer เพราะไม่ต้อง copy เร็วกว่า ดีกว่า copy struct
	return c.db
}

// implement Functions
func (d *db) Url() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}
func (d *db) MaxCon() int { return d.maxConnections }

// ------------------------------------------ JWT -----------------------------------
type IJwtConfig interface {
	SecretKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetJwtAccessExpires(t int) //เพิ่มเพื่อใช้ประกอบวันหมดอายุ
	SetJwtRefreshxpires(t int) //เพิ่มเพื่อใช้ประกอบวันหมดอายุ

}
type jwt struct {
	adminKey         string
	secretKey        string
	apiKey           string
	accessExpiresAt  int //sec
	refreshExpiresAt int //sec
}

func (c *config) Jwt() IJwtConfig { //ใช้ Pointer เพราะไม่ต้อง copy เร็วกว่า ดีกว่า copy struct
	return c.jwt
}

// implement function
func (j *jwt) SecretKey() []byte         { return []byte(j.secretKey) }
func (j *jwt) AdminKey() []byte          { return []byte(j.adminKey) }
func (j *jwt) ApiKey() []byte            { return []byte(j.apiKey) }
func (j *jwt) AccessExpiresAt() int      { return j.accessExpiresAt }
func (j *jwt) RefreshExpiresAt() int     { return j.refreshExpiresAt }
func (j *jwt) SetJwtAccessExpires(t int) { j.accessExpiresAt = t }
func (j *jwt) SetJwtRefreshxpires(t int) { j.refreshExpiresAt = t }
