package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type customLogger struct {
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

func NewLogger() *customLogger {
	return &customLogger{
		InfoLogger:  makeLogger("INFO: ", os.Stdout),
		WarnLogger:  makeLogger("WARN: ", os.Stdout),
		ErrorLogger: makeLogger("ERROR: ", os.Stdout),
	}
}

func makeLogger(levelPrefix string, file io.Writer) *log.Logger {
	return log.New(file, levelPrefix, log.Ldate|log.Ltime)
}

func (c *customLogger) Info(msg string, args ...interface{}) {
	c.InfoLogger.Println(fmt.Sprintf(msg, args...))
}

func (c *customLogger) Warn(msg string, args ...interface{}) {
	c.WarnLogger.Println(fmt.Sprintf(msg, args...))
}

func (c *customLogger) Err(msg string, args ...interface{}) {
	c.ErrorLogger.Println(fmt.Sprintf(msg, args...))
}
