package db

import (
	"context"
	"domain-admin/pkg/config"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	dbs  = make(map[string]*gorm.DB)
	lock sync.RWMutex
)

// InitDB 初始化默认数据源
func InitDB(cfg config.DataBaseConfig) {
	instance := newDB(cfg)
	setDB("default", instance)
}

// RegisterDB 注册多数据源
func RegisterDB(name string, cfg config.DataBaseConfig) {
	instance := newDB(cfg)
	setDB(name, instance)
}

// newDB 根据 Driver 初始化 gorm.DB
func newDB(cfg config.DataBaseConfig) *gorm.DB {
	dialector := getDialector(cfg)
	if dialector == nil {
		panic(fmt.Sprintf("unsupported database driver: %s", cfg.Driver))
	}

	gormConfig := &gorm.Config{
		Logger:         logger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
		NamingStrategy: namingStrategy(cfg),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatalf("failed to connect to database [%s]: %v", cfg.DB, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	return db
}

func getDialector(cfg config.DataBaseConfig) gorm.Dialector {
	dsn := cfg.DSN()
	switch strings.ToLower(cfg.Driver) {
	case "mysql":
		return mysql.Open(dsn)
	case "postgres":
		return postgres.Open(dsn)
	case "sqlite":
		return sqlite.Open(dsn)
	default:
		return nil
	}
}

func namingStrategy(cfg config.DataBaseConfig) schema.NamingStrategy {
	return schema.NamingStrategy{
		TablePrefix:   cfg.Prefix,
		SingularTable: cfg.Singular,
	}
}

func parseLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	default:
		return logger.Info
	}
}

func setDB(name string, instance *gorm.DB) {
	lock.Lock()
	defer lock.Unlock()
	dbs[name] = instance
}

func GetDB(name string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := dbs[name]
	if !ok {
		panic(fmt.Sprintf("db instance %s not found", name))
	}
	return db
}

func Default() *gorm.DB {
	return GetDB("default")
}

func WithContext(ctx context.Context, name string) *gorm.DB {
	return GetDB(name).WithContext(ctx)
}

func Transaction(name string, fc func(tx *gorm.DB) error) error {
	return GetDB(name).Transaction(fc)
}
