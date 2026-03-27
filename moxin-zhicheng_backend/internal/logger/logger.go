package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var L *slog.Logger

func InitLogger(env string) {
	var writer io.Writer
	var level slog.Level

	if env == "config_dev" {
		level = slog.LevelInfo
		fileWriter := &lumberjack.Logger{
			Filename: "logs/app.log",
			MaxSize:  100,
			Compress: true,
		}
		// 同时输出到控制台和文件
		writer = io.MultiWriter(os.Stdout, fileWriter)
	} else {
		level = slog.LevelDebug
		writer = os.Stdout
	}

	// 这里的 NewJSONHandler 内部其实已经处理了基本的并发安全
	L = slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // 打印行号
	}))

	slog.SetDefault(L)
}

// --- 以下是二次封装的方法 ---

// Info 封装，支持 kv 对
func Info(msg string, args ...any) {
	L.Info(msg, args...)
}

// Error 封装，自动记录错误信息
func Error(msg string, err error, args ...any) {
	// 将 error 放到属性中，而不是拼接到 msg 字符串里，方便日志分析
	newArgs := append([]any{slog.Any("error", err)}, args...)
	L.Error(msg, newArgs...)
}

// Debug 封装
func Debug(msg string, args ...any) {
	L.Debug(msg, args...)
}

// Warn 封装
func Warn(msg string, args ...any) {
	L.Warn(msg, args...)
}

// WithCtx 这是一个高级技巧：如果你以后加了链路追踪，可以用这个方法把 trace_id 带入
func WithCtx(ctx context.Context) *slog.Logger {
	// 假设你从 context 里取出了 trace_id
	// traceID := ctx.Value("trace_id").(string)
	// return L.With(slog.String("trace_id", traceID))
	return L
}
