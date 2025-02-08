package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"todo/pkg/config"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

// Init 初始化日志
func Init(cfg config.LoggerConfig) error {
	// 设置日志级别
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("无效的日志级别: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	// 配置输出
	output := configureOutput(cfg.File)

	// 设置日志格式
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error"
	zerolog.TimeFieldFormat = time.RFC3339

	// 自定义日志级别显示
	zerolog.LevelDebugValue = "调试"
	zerolog.LevelInfoValue = "信息"
	zerolog.LevelWarnValue = "警告"
	zerolog.LevelErrorValue = "错误"
	zerolog.LevelFatalValue = "致命"

	// 初始化日志对象
	log = zerolog.New(output).With().Timestamp().Logger().Hook(TimezoneHook{})

	return nil
}

func configureOutput(filePath string) io.Writer {
	if filePath == "" {
		return os.Stdout
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Printf("创建日志目录失败: %v, 使用标准输出", err)
		return os.Stdout
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("打开日志文件失败: %v, 使用标准输出", err)
		return os.Stdout
	}

	return file
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

// LogRequest 记录HTTP请求日志
func LogRequest(traceID string, method string, path string, status int, latency time.Duration) {
	log.Info().
		Str("追踪ID", traceID).
		Str("请求方法", method).
		Str("请求路径", path).
		Int("状态码", status).
		Dur("响应时间", latency).
		Msg("HTTP请求")
}

// 添加一些常用的日志方法
func LogDBOperation(operation string, table string, err error) {
	if err != nil {
		log.Error().
			Str("操作", operation).
			Str("数据表", table).
			Err(err).
			Msg("数据库操作失败")
	} else {
		log.Debug().
			Str("操作", operation).
			Str("数据表", table).
			Msg("数据库操作成功")
	}
}

func LogCacheOperation(operation string, key string, err error) {
	if err != nil {
		log.Error().
			Str("操作", operation).
			Str("键", key).
			Err(err).
			Msg("缓存操作失败")
	} else {
		log.Debug().
			Str("操作", operation).
			Str("键", key).
			Msg("缓存操作成功")
	}
}

func LogUserAction(userID string, action string, details string) {
	log.Info().
		Str("用户ID", userID).
		Str("操作", action).
		Str("详情", details).
		Msg("用户活动")
}

func LogSystemEvent(event string, details string) {
	log.Info().
		Str("事件", event).
		Str("详情", details).
		Msg("系统事件")
}

// TimezoneHook 用于确保时间戳使用正确的时区
type TimezoneHook struct{}

func (h TimezoneHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Time("时间", time.Now().In(time.Local))
}
