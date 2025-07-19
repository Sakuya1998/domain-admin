package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig          `mapstructure:"server"`
	Database      DataBaseConfig        `mapstructure:"database"`
	Redis         RedisConfig           `mapstructure:"redis"`
	Log           LogConfig             `mapstructure:"log"`
	JWT           JWTConfig             `mapstructure:"jwt"`
	CloudProvider []CloudProviderConfig `mapstructure:"cloudprovider"`
	OTLP          OTLPConfig            `mapstructure:"otel"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	Timeout string `mapstructure:"timeout"`
}

type DataBaseConfig struct {
	Driver   string `mapstructure:"driver"`    // mysql, postgres, sqlite
	Host     string `mapstructure:"host"`      // 主机
	Port     int    `mapstructure:"port"`      // 端口
	User     string `mapstructure:"user"`      // 用户名
	Pass     string `mapstructure:"password"`  // 密码
	DB       string `mapstructure:"db"`        // 数据库名
	Charset  string `mapstructure:"charset"`   // 字符集，默认 utf8mb4
	SSLMode  string `mapstructure:"ssl_mode"`  // Postgres 用
	FilePath string `mapstructure:"file_path"` // SQLite 用
	Prefix   string `mapstructure:"prefix"`    // 表名前缀
	Singular bool   `mapstructure:"singular"`  // 是否使用单数表名
	Engine   string `mapstructure:"engine"`    // MySQL 表引擎，如 InnoDB

	MaxIdleConns    int `mapstructure:"max_idle_conns"`    // 连接池空闲数
	MaxOpenConns    int `mapstructure:"max_open_conns"`    // 连接池最大连接数
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"` // 单位分钟

	LogLevel string `mapstructure:"log_level"` // 日志级别：silent, error, warn, info
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

type LogConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	Output      string `mapstructure:"output"`
	FilePath    string `mapstructure:"file_path"`
	MaxSize     int    `mapstructure:"max_size"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"`
	Compress    bool   `mapstructure:"compress"`
	Development bool   `mapstructure:"development"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration string `mapstructure:"expiration"`
}

type CloudProviderConfig struct {
	Type         string `mapstructure:"type"`
	AccessKey    string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
}

type OTLPConfig struct {
	Endpoint          string `mapstructure:"endpoint"`
	ServiceName       string `mapstructure:"service_name"`
	ServiceVersion    string `mapstructure:"service_version"`
	ServiceInstanceID string `mapstructure:"service_instance_id"`
}

var cfg = &Config{}

func InitConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Config unmarshal error: %v", err)
	}

	fmt.Println("[Config] loaded successfully")
	return cfg
}

func GetConfig() *Config {
	return cfg
}

func (c DataBaseConfig) DSN() string {
	switch c.Driver {
	case "mysql":
		charset := c.Charset
		if charset == "" {
			charset = "utf8mb4"
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.User, c.Pass, c.Host, c.Port, c.DB, charset)
	case "postgres":
		sslMode := c.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			c.Host, c.Port, c.User, c.DB, c.Pass, sslMode)
	case "sqlite":
		if c.FilePath != "" {
			return c.FilePath
		}
		return c.DB + ".db"
	default:
		return ""
	}
}
