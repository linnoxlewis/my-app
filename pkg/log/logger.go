package log

import (
	"fmt"
	"log"
)

type Logging interface {
	LogErrorf(format string, v ...any)
	LogInfo(format string, v ...any)
}

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{log.Default()}
}

func (l *Logger) LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("[Error]: %s\n", msg)
}

func (l *Logger) LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.Printf("[Info]: %s\n", msg)
}
