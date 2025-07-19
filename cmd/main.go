package main

import (
	"domain-admin/api"
	"domain-admin/internal/migration"
	"domain-admin/pkg/cache"
	"domain-admin/pkg/config"
	"domain-admin/pkg/db"
	"domain-admin/pkg/logger"
	"domain-admin/pkg/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	cfg := config.GetConfig()

	logger.InitLogger(cfg.Log)
	db.InitDB(cfg.Database)
	cache.InitCache(cfg.Redis)

	// 数据库迁移
	if err := migration.AutoMigrate(db.GetDB("default")); err != nil {
		logger.Errorf("数据库迁移失败: %v", err)
		panic(err)
	}

	// 创建默认管理员
	if err := migration.CreateDefaultAdmin(db.GetDB("default")); err != nil {
		logger.Errorf("创建默认管理员失败: %v", err)
		panic(err)
	}

	// 初始化RBAC基础数据
	if err := migration.InitRBACData(db.GetDB("default")); err != nil {
		logger.Errorf("初始化RBAC基础数据失败: %v", err)
		panic(err)
	}

	// 初始化RBAC系统
	if err := middleware.InitRBAC(db.GetDB("default"), "configs/rbac_model.conf"); err != nil {
		logger.Errorf("初始化RBAC系统失败: %v", err)
		panic(err)
	}

	// 同步RBAC策略
	if err := middleware.SyncRBACPolicies(db.GetDB("default")); err != nil {
		logger.Errorf("同步RBAC策略失败: %v", err)
		panic(err)
	}

	r := gin.Default()
	api.RegisterRoutes(r)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}
