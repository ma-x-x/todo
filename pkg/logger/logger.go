package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"todo/pkg/config"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

// LogContext 日志上下文
type LogContext struct {
	TraceID    string
	RequestID  string
	UserID     uint
	Additional map[string]interface{}
}

// TimezoneHook 时区钩子
type TimezoneHook struct{}

func (h TimezoneHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Time("timestamp", time.Now())
}

// Init 初始化日志
func Init(cfg config.LoggerConfig) error {
	// 配置日志轮转
	output := configureOutput(cfg.File)

	// 设置日志级别
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("无效的日志级别: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	// 设置日志格式
	configureLogFormat()

	// 初始化日志对象
	log = zerolog.New(output).With().Timestamp().Logger().Hook(TimezoneHook{})
	return nil
}

// configureLogFormat 配置日志格式
func configureLogFormat() {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error"
	zerolog.TimeFieldFormat = time.RFC3339

	zerolog.LevelDebugValue = "debug"
	zerolog.LevelInfoValue = "info"
	zerolog.LevelWarnValue = "warn"
	zerolog.LevelErrorValue = "error"
	zerolog.LevelFatalValue = "fatal"
}

func configureOutput(filePath string) io.Writer {
	if filePath == "" {
		return os.Stdout
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v, using stdout\n", err)
		return os.Stdout
	}

	logWriter := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	return zerolog.MultiLevelWriter(os.Stdout, logWriter)
}

// 基础日志方法
func Debug() *zerolog.Event { return log.Debug() }
func Info() *zerolog.Event  { return log.Info() }
func Warn() *zerolog.Event  { return log.Warn() }
func Error() *zerolog.Event { return log.Error() }
func Fatal() *zerolog.Event { return log.Fatal() }

// WithContext 创建带上下文的日志事件
func WithContext(ctx LogContext) zerolog.Logger {
	contextLogger := log.With().
		Str("trace_id", ctx.TraceID).
		Uint("user_id", ctx.UserID).
		Time("timestamp", time.Now())

	if ctx.RequestID != "" {
		contextLogger = contextLogger.Str("request_id", ctx.RequestID)
	}

	for k, v := range ctx.Additional {
		contextLogger = contextLogger.Interface(k, v)
	}

	return contextLogger.Logger()
}

// LogOperation 日志操作结构
type LogOperation struct {
	TraceID    string
	Operation  string
	Duration   time.Duration
	Error      error
	Additional map[string]interface{}
}

// LogHTTPRequest 记录HTTP请求
func LogHTTPRequest(op LogOperation, method, path string, status int) {
	event := Info().
		Str("trace_id", op.TraceID).
		Str("method", method).
		Str("path", path).
		Int("status", status).
		Dur("duration", op.Duration)

	if op.Error != nil {
		event.Err(op.Error)
	}

	for k, v := range op.Additional {
		event.Interface(k, v)
	}

	event.Msg("HTTP Request")
}

// LogDBOperation 记录数据库操作
func LogDBOperation(op LogOperation, table string) {
	event := Debug().
		Str("trace_id", op.TraceID).
		Str("operation", op.Operation).
		Str("table", table).
		Dur("duration", op.Duration)

	if op.Error != nil {
		event.Err(op.Error).Msg("Database operation failed")
	} else {
		event.Msg("Database operation succeeded")
	}
}

// LogCacheOperation 记录缓存操作
func LogCacheOperation(op LogOperation, key string) {
	event := Debug().
		Str("trace_id", op.TraceID).
		Str("operation", op.Operation).
		Str("key", key).
		Dur("duration", op.Duration)

	if op.Error != nil {
		event.Err(op.Error)
	}

	event.Msg("Cache Operation")
}

// LogInfo 记录信息日志
func LogInfo(traceID string, message string, fields map[string]interface{}) {
	event := log.Info().
		Str("trace_id", traceID)

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(message)
}

// LogError 记录错误日志
func LogError(traceID string, err error, message string, fields map[string]interface{}) {
	event := log.Error().
		Str("trace_id", traceID)

	if err != nil {
		event = event.Err(err)
	}

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(message)
}

// LogDebug 记录调试日志
func LogDebug(traceID string, msg string, fields map[string]interface{}) {
	event := log.Debug().
		Str("trace_id", traceID)

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}

// LogWarn 记录警告日志
func LogWarn(traceID string, msg string, fields map[string]interface{}) {
	event := log.Warn().
		Str("trace_id", traceID)

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}

// LogDBQuery 记录数据库查询日志
func LogDBQuery(traceID string, query string, args []interface{}, duration time.Duration) {
	log.Debug().
		Str("trace_id", traceID).
		Str("query", query).
		Interface("args", args).
		Dur("duration", duration).
		Msg("Database Query")
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

// 1. 添加结构化日志字段
type LogField struct {
	Key   string
	Value interface{}
}

// 2. 支持日志轮转
func InitLogger(cfg config.LoggerConfig) error {
	// 配置日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.File,
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	}

	// 多输出支持
	var output io.Writer
	if cfg.File != "" {
		output = zerolog.MultiLevelWriter(os.Stdout, lumberjackLogger)
	} else {
		output = os.Stdout
	}

	// 设置日志级别
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("无效的日志级别: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	// 初始化日志对象
	log = zerolog.New(output).
		With().
		Timestamp().
		Logger().
		Hook(TimezoneHook{})

	return nil
}
