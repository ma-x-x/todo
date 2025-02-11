package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`              // 数据库主机地址
	Port            int           `mapstructure:"port"`              // 数据库端口
	Username        string        `mapstructure:"username"`          // 数据库用户名
	Password        string        `mapstructure:"password"`          // 数据库密码
	Database        string        `mapstructure:"database"`          // 数据库名称
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
	Slaves          []string      `mapstructure:"slaves"`            // 新增：从库配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode         string `mapstructure:"mode"`          // 服务器模式（debug/release）
	Port         int    `mapstructure:"port"`          // 服务器端口
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 读取超时时间（秒）
	WriteTimeout int    `mapstructure:"write_timeout"` // 写入超时时间（秒）
	IdleTimeout  int    `mapstructure:"idle_timeout"`  // 空闲超时时间（秒）
	SwaggerHost  string `mapstructure:"swagger_host"`  // 添加这行
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`      // Redis主机地址
	Port     int    `mapstructure:"port"`      // Redis端口
	Password string `mapstructure:"password"`  // Redis密码
	DB       int    `mapstructure:"db"`        // Redis数据库索引
	PoolSize int    `mapstructure:"pool_size"` // 连接池大小
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level string `mapstructure:"level"` // 日志级别（debug/info/warn/error）
	File  string `mapstructure:"file"`  // 日志文件路径
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`       // JWT密钥
	ExpireHours int    `mapstructure:"expire_hours"` // JWT过期时间（小时）
	Issuer      string `mapstructure:"issuer"`       // JWT签发者
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

// Config 应用配置
// 配置加载优先级（从高到低）：
// 1. 环境变量（例如：DB_HOST, REDIS_PORT）
// 2. 配置文件（通过 CONFIG_FILE 环境变量指定的文件）
// 3. 默认配置文件（基于 APP_ENV 环境变量，如 config.prod.yaml）
// 4. 代码中的默认值
type Config struct {
	Server    ServerConfig    `mapstructure:"server" validate:"required"`
	MySQL     MySQLConfig     `mapstructure:"mysql" validate:"required"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Logger    LoggerConfig    `mapstructure:"logger"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	TaskQueue struct {
		BufferSize int `mapstructure:"buffer_size"` // 任务队列缓冲大小
		Workers    int `mapstructure:"workers"`     // 工作协程数量
	} `mapstructure:"task_queue"`
}

// LoadConfig 改名为 Load
func Load() (*Config, error) {
	// 设置配置文件路径
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "configs/config.yaml"
	}

	// 设置默认值和环境变量绑定
	setDefaults()
	bindEnvVariables()

	// 设置并读取配置文件
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	log.Printf("配置加载成功 [模式:%s] [端口:%d] [日志级别:%s]",
		cfg.Server.Mode, cfg.Server.Port, cfg.Logger.Level)

	return &cfg, nil
}

// setDefaults 设置默认配置值
// 这些默认值的优先级最低，会被配置文件和环境变量覆盖
func setDefaults() {
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.idle_timeout", 120)

	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)
	viper.SetDefault("mysql.conn_max_lifetime", "1h")

	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.file", "logs/app.log")

	viper.SetDefault("jwt.expire_hours", 1)
	viper.SetDefault("jwt.issuer", "todo_app")
}

// processEnvVars 处理环境变量替换
// 支持的环境变量及其对应的配置项：
// - DB_HOST: mysql.host
// - DB_PORT: mysql.port
// - DB_USER: mysql.username
// - DB_PASSWORD: mysql.password
// - DB_NAME: mysql.database
// - REDIS_HOST: redis.host
// - REDIS_PORT: redis.port
// - REDIS_PASSWORD: redis.password
// - JWT_SECRET: jwt.secret
// - LOG_LEVEL: logger.level
func processEnvVars() {
	fmt.Println("========== 开始处理环境变量 ==========")

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		fmt.Printf("从环境变量读取 DB_HOST: %s\n", host)
		viper.Set("mysql.host", host)
	} else {
		fmt.Println("未找到环境变量 DB_HOST")
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		fmt.Printf("从环境变量读取 DB_PORT: %s\n", port)
		viper.Set("mysql.port", port)
	} else {
		fmt.Println("未找到环境变量 DB_PORT")
	}

	if user := os.Getenv("DB_USER"); user != "" {
		fmt.Printf("从环境变量读取 DB_USER: %s\n", user)
		viper.Set("mysql.username", user)
	} else {
		fmt.Println("未找到环境变量 DB_USER")
	}

	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		fmt.Println("从环境变量读取 DB_PASSWORD: ******")
		viper.Set("mysql.password", pass)
	} else {
		fmt.Println("未找到环境变量 DB_PASSWORD")
	}

	if name := os.Getenv("DB_NAME"); name != "" {
		fmt.Printf("从环境变量读取 DB_NAME: %s\n", name)
		viper.Set("mysql.database", name)
	} else {
		fmt.Println("未找到环境变量 DB_NAME")
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		fmt.Printf("从环境变量读取 REDIS_HOST: %s\n", host)
		viper.Set("redis.host", host)
	} else {
		fmt.Println("未找到环境变量 REDIS_HOST")
	}

	if port := os.Getenv("REDIS_PORT"); port != "" {
		fmt.Printf("从环境变量读取 REDIS_PORT: %s\n", port)
		viper.Set("redis.port", port)
	} else {
		fmt.Println("未找到环境变量 REDIS_PORT")
	}

	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		fmt.Println("从环境变量读取 REDIS_PASSWORD: ******")
		viper.Set("redis.password", pass)
	} else {
		fmt.Println("未找到环境变量 REDIS_PASSWORD")
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		fmt.Println("从环境变量读取 JWT_SECRET: ******")
		viper.Set("jwt.secret", secret)
	} else {
		fmt.Println("未找到环境变量 JWT_SECRET")
	}

	// 日志配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		fmt.Printf("从环境变量读取 LOG_LEVEL: %s\n", level)
		viper.Set("logger.level", level)
	} else {
		fmt.Println("未找到环境变量 LOG_LEVEL")
	}

	fmt.Println("========== 环境变量处理完成 ==========")
}

// bindEnvVariables 绑定环境变量
// 允许使用环境变量覆盖任何配置项
// 环境变量名称规则：配置路径中的点号(.)替换为下划线(_)
// 例如：mysql.host 对应环境变量 MYSQL_HOST
func bindEnvVariables() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// validateConfig 验证配置
// 检查必要的配置项并设置默认值
// 必要的环境变量：
// - DB_PASSWORD: 数据库密码
// - JWT_SECRET: JWT密钥
// 默认值：
// - MySQL Host: mysql
// - MySQL Port: 3306
// - MySQL Username: todo_user
// - MySQL Database: todo_db
func validateConfig(cfg *Config) error {
	var missingVars []string

	if cfg.MySQL.Password == "" {
		missingVars = append(missingVars, "DB_PASSWORD")
	}
	if cfg.JWT.Secret == "" {
		missingVars = append(missingVars, "JWT_SECRET")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("缺少必要的环境变量: %v", missingVars)
	}

	// 验证数据库配置
	if cfg.MySQL.Host == "" {
		cfg.MySQL.Host = "mysql" // 默认值
	}
	if cfg.MySQL.Port == 0 {
		cfg.MySQL.Port = 3306 // 默认值
	}
	if cfg.MySQL.Username == "" {
		cfg.MySQL.Username = "todo_user" // 默认值
	}
	if cfg.MySQL.Database == "" {
		cfg.MySQL.Database = "todo_db" // 默认值
	}

	// 验证服务器模式
	if cfg.Server.Mode != "debug" && cfg.Server.Mode != "release" && cfg.Server.Mode != "test" {
		log.Printf("警告：配置文件中的服务器模式 '%s' 无效，使用默认的 'release' 模式", cfg.Server.Mode)
		cfg.Server.Mode = "release"
	}

	return nil
}

// Validate 使用结构化的配置验证
func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// WatchConfig 添加配置热重载支持
func WatchConfig(cfg *Config, callback func(cfg *Config)) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 重新加载配置
	})
}
