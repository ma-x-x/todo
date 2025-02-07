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
func bindEnvVariables() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// validateConfig 验证配置
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
