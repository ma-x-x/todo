package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`             // 数据库主机地址
	Port            int           `mapstructure:"port"`             // 数据库端口
	Username        string        `mapstructure:"username"`         // 数据库用户名
	Password        string        `mapstructure:"password"`         // 数据库密码
	Database        string        `mapstructure:"database"`         // 数据库名称
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`   // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns"`   // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode         string `mapstructure:"mode"`          // 服务器模式（debug/release）
	Port         int    `mapstructure:"port"`          // 服务器端口
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 读取超时时间（秒）
	WriteTimeout int    `mapstructure:"write_timeout"` // 写入超时时间（秒）
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

// Config 应用配置
// 配置加载优先级（从高到低）：
// 1. 环境变量（例如：DB_HOST, REDIS_PORT）
// 2. 配置文件（通过 CONFIG_FILE 环境变量指定的文件）
// 3. 默认配置文件（基于 APP_ENV 环境变量，如 config.prod.yaml）
// 4. 代码中的默认值
type Config struct {
	Server    ServerConfig `mapstructure:"server"`
	MySQL     MySQLConfig  `mapstructure:"mysql"`
	Redis     RedisConfig  `mapstructure:"redis"`
	Logger    LoggerConfig `mapstructure:"logger"`
	JWT       JWTConfig    `mapstructure:"jwt"`
	RateLimit struct {
		RequestsPerSecond float64 `mapstructure:"requests_per_second"` // 每秒请求限制
		Burst             int     `mapstructure:"burst"`               // 突发请求限制
	} `mapstructure:"rate_limit"`
	TaskQueue struct {
		BufferSize int `mapstructure:"buffer_size"` // 任务队列缓冲大小
		Workers    int `mapstructure:"workers"`     // 工作协程数量
	} `mapstructure:"task_queue"`
}

// LoadConfig 加载配置文件
// 配置加载流程：
// 1. 设置默认值
// 2. 绑定环境变量
// 3. 读取配置文件
// 4. 处理环境变量覆盖
// 5. 解析配置到结构体
// 6. 验证必要的配置
func LoadConfig() (*Config, error) {
	// 1. 设置默认值
	setDefaults()

	// 2. 绑定环境变量
	bindEnvVariables()

	// 3. 读取配置文件
	if err := loadConfigFile(); err != nil {
		return nil, err
	}

	// 4. 处理环境变量替换
	processEnvVars()

	// 5. 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 6. 验证必要的配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults 设置默认配置值
// 这些默认值的优先级最低，会被配置文件和环境变量覆盖
func setDefaults() {
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)

	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)
	viper.SetDefault("mysql.conn_max_lifetime", "1h")

	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.file", "logs/app.log")
}

// loadConfigFile 加载配置文件
// 配置文件加载优先级：
// 1. CONFIG_FILE 环境变量指定的文件
// 2. 基于 APP_ENV 的配置文件（如 config.prod.yaml）
// 3. 在以下路径查找：./configs、../configs、/app/configs
func loadConfigFile() error {
	viper.SetConfigType("yaml")

	// 获取配置文件路径
	configFile := os.Getenv("CONFIG_FILE")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// 根据环境选择配置文件
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "dev" // 默认开发环境
		}

		configName := fmt.Sprintf("config.%s.yaml", env)
		viper.SetConfigName(configName)
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")
		viper.AddConfigPath("/app/configs")
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
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
	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		viper.Set("mysql.host", host)
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		viper.Set("mysql.port", port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		viper.Set("mysql.username", user)
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		viper.Set("mysql.password", pass)
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		viper.Set("mysql.database", name)
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		viper.Set("redis.host", host)
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		viper.Set("redis.port", port)
	}
	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		viper.Set("redis.password", pass)
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		viper.Set("jwt.secret", secret)
	}

	// 日志配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		viper.Set("logger.level", level)
	}
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
		return fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	// 验证数据库配置
	if cfg.MySQL.Host == "" {
		cfg.MySQL.Host = "mysql"  // 默认值
	}
	if cfg.MySQL.Port == 0 {
		cfg.MySQL.Port = 3306  // 默认值
	}
	if cfg.MySQL.Username == "" {
		cfg.MySQL.Username = "todo_user"  // 默认值
	}
	if cfg.MySQL.Database == "" {
		cfg.MySQL.Database = "todo_db"  // 默认值
	}

	return nil
}
