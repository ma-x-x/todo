package logger

import (
	"fmt"
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
	fmt.Errorf("日志级别: %v\n", level)
	if err != nil {
		return fmt.Errorf("无效的日志级别: %v", err)
	}
	zerolog.SetGlobalLevel(level)

	// 设置输出
	var output *os.File
	if cfg.File != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %v", err)
		}

		// 打开日志文件
		output, err = os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("打开日志文件失败: %v", err)
		}
	} else {
		output = os.Stdout
	}

	// 设置日志格式
	zerolog.TimestampFieldName = "时间"
	zerolog.LevelFieldName = "级别"
	zerolog.MessageFieldName = "消息"
	zerolog.ErrorFieldName = "错误"

	// 设置时间格式
	zerolog.TimeFieldFormat = time.RFC3339Nano

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
