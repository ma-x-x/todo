package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"todo-demo/pkg/config"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

// Init 初始化日志
func Init(cfg config.LoggerConfig) error {
	// 设置日志级别
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	zerolog.SetGlobalLevel(level)

	// 设置输出
	var output *os.File
	if cfg.File != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}

		// 打开日志文件
		output, err = os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
	} else {
		output = os.Stdout
	}

	// 初始化日志对象
	log = zerolog.New(output).With().Timestamp().Logger()

	return nil
}

// Debug 输出调试日志
func Debug() *zerolog.Event {
	return log.Debug()
}

// Info 输出信息日志
func Info() *zerolog.Event {
	return log.Info()
}

// Warn 输出警告日志
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error 输出错误日志
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal 输出致命错误日志
func Fatal() *zerolog.Event {
	return log.Fatal()
}

func LogRequest(traceID string, method string, path string, status int, latency time.Duration) {
	log.Info().
		Str("trace_id", traceID).
		Str("method", method).
		Str("path", path).
		Int("status", status).
		Dur("latency", latency).
		Msg("HTTP Request")
}
