package main

import (
	"fmt"
	"os"

	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/servers"
	"github.com/PHURINTOR/phurinshop/pkg/database"
)

// ----------------------------------------- Function path Run main.exe -------------------------
func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1] //เวลารัน main.exe dev   = dev คือ อากิวเม้น ที่ อยู่ในลำดับ 2
	}
}

func main() {
	fmt.Println("Hello")
	cfg := config.LoadConfig(envPath())
	//Run = air -c .air.dev.toml

	//Connect Db
	db := database.DbConnect(cfg.Db())
	defer db.Close()

	fmt.Print(db)

	//Start server Gofiber
	servers.NewServer(cfg, db).Start()
}
