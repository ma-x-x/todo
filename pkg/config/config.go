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
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
	Issuer      string `mapstructure:"issuer"`
}

// Config 应用配置
type Config struct {
	Server    ServerConfig `mapstructure:"server"`
	MySQL     MySQLConfig  `mapstructure:"mysql"`
	Redis     RedisConfig  `mapstructure:"redis"`
	Logger    LoggerConfig `mapstructure:"logger"`
	JWT       JWTConfig    `mapstructure:"jwt"`
	RateLimit struct {
		RequestsPerSecond float64 `mapstructure:"requests_per_second"`
		Burst             int     `mapstructure:"burst"`
	} `mapstructure:"rate_limit"`
	TaskQueue struct {
		BufferSize int `mapstructure:"buffer_size"`
		Workers    int `mapstructure:"workers"`
	} `mapstructure:"task_queue"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 1. 设置默认值
	setDefaults()

	// 2. 读取配置文件
	if err := loadConfigFile(); err != nil {
		return nil, err
	}

	// 3. 绑定环境变量
	bindEnvVariables()

	// 4. 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 5. 验证必要的配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults 设置默认配置值
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

		configName := fmt.Sprintf("config.%s.yaml", env) // 添加.yaml后缀
		viper.SetConfigName(configName)
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")  // 添加上级目录
		viper.AddConfigPath("/app/configs")
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果是开发环境，尝试使用默认配置文件
		if os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "dev" {
			viper.SetConfigName("config")  // 尝试读取无环境后缀的配置
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

// bindEnvVariables 绑定环境变量
func bindEnvVariables() {
	// 自动绑定环境变量，环境变量格式为：APP_SERVER_PORT
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 显式绑定关键环境变量
	viper.BindEnv("mysql.password", "DB_PASSWORD")
	viper.BindEnv("mysql.host", "DB_HOST")
	viper.BindEnv("mysql.port", "DB_PORT")
	viper.BindEnv("mysql.username", "DB_USER")
	viper.BindEnv("mysql.database", "DB_NAME")
	
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")
}

// validateConfig 验证配置
func validateConfig(cfg *Config) error {
	// 验证必要的配置项
	if cfg.MySQL.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	return nil
}
