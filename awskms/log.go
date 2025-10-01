package awskms

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/smithy-go/logging"
	"github.com/yyle88/zaplog"
)

// SlogLogger implements logging.Logger interface using standard slog package
// Provides basic logging function in AWS SDK operations with formatted message output
// Outputs each message at Info grade with classification prefix in operation tracking
//
// SlogLogger 使用标准 slog 包实现 logging.Logger 接口
// 在 AWS SDK 操作中提供基本日志记录功能，带有格式化消息输出
// 以 Info 级别输出每条消息，带有分类前缀以跟踪操作
type SlogLogger struct{}

// NewSlogLogger creates a new SlogLogger instance in AWS SDK logging
// Returns zero-value instance that uses standard slog package in output
// Suitable in applications that choose standard slog solution
//
// NewSlogLogger 在 AWS SDK 日志记录中创建新的 SlogLogger 实例
// 返回使用标准 slog 包输出的零值实例
// 适用于选择标准 slog 解决方案的应用
func NewSlogLogger() *SlogLogger {
	return &SlogLogger{}
}

// Logf implements logging.Logger interface using slog with formatted message output
// Formats message with classification prefix and outputs using standard slog package
// Provides consistent logging format in AWS SDK operations tracking
//
// Logf 使用 slog 实现 logging.Logger 接口，带有格式化消息输出
// 使用分类前缀格式化消息并使用标准 slog 包输出
// 在 AWS SDK 操作跟踪中提供一致的日志格式
func (L *SlogLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	slog.InfoContext(context.Background(), fmt.Sprintf("[%s] %s", classification, fmt.Sprintf(format, v...)))
}

// ZapLogger implements logging.Logger interface using uber/zap
// Provides structured logging function in AWS SDK operations with accurate stack skip
// Uses zaplog.ZAPS.Skip(1) to adjust the stack in accurate source location reporting
//
// ZapLogger 使用 uber/zap 实现 logging.Logger 接口
// 在 AWS SDK 操作中提供结构化日志记录功能，带有准确的栈跳过
// 使用 zaplog.ZAPS.Skip(1) 调整栈以准确报告源位置
type ZapLogger struct{}

// NewZapLogger creates a new ZapLogger instance in AWS SDK logging
// Returns zero-value instance that uses uber/zap in structured output
// Suitable in applications requiring advanced structured logging capabilities
//
// NewZapLogger 在 AWS SDK 日志记录中创建新的 ZapLogger 实例
// 返回使用 uber/zap 进行结构化输出的零值实例
// 适用于需要高级结构化日志记录能力的应用
func NewZapLogger() *ZapLogger {
	return &ZapLogger{}
}

// Logf implements logging.Logger interface using zap with stack skip adjustment
// Formats message with classification prefix and outputs using uber/zap structured mechanism
// Adjusts the stack with Skip(1) to ensure correct source location in output
//
// Logf 使用 zap 实现 logging.Logger 接口，带有栈跳过调整
// 使用分类前缀格式化消息并使用 uber/zap 结构化机制输出
// 使用 Skip(1) 调整栈以确保输出中的源位置正确
func (L *ZapLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	zaplog.ZAPS.Skip(1).SUG.Infof("[%s] %s", classification, fmt.Sprintf(format, v...))
}
