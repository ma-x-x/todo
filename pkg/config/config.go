package config

import (
	"fmt"
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
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// 添加调试日志
	fmt.Printf("Looking for config file in: %s\n", viper.ConfigFileUsed())

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 打印数据库配置
	fmt.Printf("MySQL Config: %+v\n", config.MySQL)

	return &config, nil
}
