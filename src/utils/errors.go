package utils

import (
	"fmt"
	"runtime"
)

func Error(err error) error {
	_, filename, line, _ := runtime.Caller(1)
	return fmt.Errorf("ERROR %s:%d %v", filename, line, err)
}
