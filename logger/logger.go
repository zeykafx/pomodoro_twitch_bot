package logger

import (
	"fmt"
	"time"
)

func Log(formattedString string) {
	fmt.Printf("[BOT-LOG] %s | %s\n", formattedString, time.Now().Local())
}
