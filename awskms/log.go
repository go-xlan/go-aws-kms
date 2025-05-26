package awskms

import (
	"fmt"
	"log/slog"
)

type SlogLogger struct{}

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{}
}

func (L *SlogLogger) Log(args ...interface{}) {
	slog.Info(fmt.Sprintln(args...))
}
