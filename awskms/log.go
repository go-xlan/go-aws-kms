package awskms

import (
	"fmt"
	"log/slog"

	"github.com/yyle88/zaplog"
)

type SlogLogger struct{}

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{}
}

func (L *SlogLogger) Log(args ...interface{}) {
	slog.Info(fmt.Sprintln(args...))
}

type ZapLogger struct{}

func NewZapLogger() *ZapLogger {
	return &ZapLogger{}
}

func (L *ZapLogger) Log(args ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Infoln(args...)
}
