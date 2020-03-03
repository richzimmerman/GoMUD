package utils

import (
	"fmt"
	"runtime"
)

func Error(err error) error {
	/*
	Custom error handler to throw a stack trace for investigation
	 */
	if err == nil {
		return nil
	}
	_, filename, line, _ := runtime.Caller(1)
	return fmt.Errorf("ERROR %s:%d %v\n", filename, line, err)
}
