package main

import (
	"fmt"
	"os"

	"github.com/PHURINTOR/phurinshop/config"
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
	fmt.Println(cfg.App())
	fmt.Println(cfg.Db())
	fmt.Println(cfg.Jwt())

	//Run = air -c .air.dev.toml
}
