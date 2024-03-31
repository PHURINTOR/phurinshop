package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Gen Name Of file
func RanFileName(ext string) string {
	filename := fmt.Sprintf("%s_%v", strings.ReplaceAll(uuid.NewString()[:6], "-", ""), time.Now().UnixMilli())
	if ext != "" {
		filename += fmt.Sprintf(".%s", ext)
	}
	return filename
}
