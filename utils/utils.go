package utils

import (
	"fmt"
	"os"
)

func Fatal(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)

	os.Exit(1)
}
