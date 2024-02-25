package utils

import (
	"encoding/json"
	"fmt"
)

func Debug(data any) {
	bytes, _ := json.MarshalIndent(data, "", "\t") //"" prefix ไม่ต้อง ,  \t คือจััดให้สวย  = ค่าจะได้เป็น ฺ byte
	fmt.Println(string(bytes))
}

func Output(data any) []byte {
	bytes, _ := json.Marshal(data) //data เพียว
	return bytes
}
